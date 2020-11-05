package job

import (
	"context"
	"sync"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

//go:generate mockery --name ORM --output ./mocks/ --case=underscore

type ORM interface {
	ListenForNewJobs() (postgres.Subscription, error)
	ClaimUnclaimedJobs(ctx context.Context) ([]models.JobSpecV2, error)
	CreateJob(ctx context.Context, jobSpec *models.JobSpecV2, taskDAG pipeline.TaskDAG) error
	DeleteJob(ctx context.Context, id int32) error
	Close() error
}

type orm struct {
	db                  *gorm.DB
	config              Config
	advisoryLocker      postgres.AdvisoryLocker
	advisoryLockClassID int32
	pipelineORM         pipeline.ORM
	eventBroadcaster    postgres.EventBroadcaster
	claimedJobs         []models.JobSpecV2
	claimedJobsMu       *sync.Mutex
}

var _ ORM = (*orm)(nil)

func NewORM(db *gorm.DB, config Config, pipelineORM pipeline.ORM, eventBroadcaster postgres.EventBroadcaster, advisoryLocker postgres.AdvisoryLocker) *orm {
	return &orm{
		db:                  db,
		config:              config,
		advisoryLocker:      advisoryLocker,
		advisoryLockClassID: postgres.AdvisoryLockClassID_JobSpawner,
		pipelineORM:         pipelineORM,
		eventBroadcaster:    eventBroadcaster,
		claimedJobs:         make([]models.JobSpecV2, 0),
		claimedJobsMu:       &sync.Mutex{},
	}
}

func (o *orm) Close() error {
	return nil
}

func (o *orm) ListenForNewJobs() (postgres.Subscription, error) {
	return o.eventBroadcaster.Subscribe(postgres.ChannelJobCreated, "")
}

// ClaimUnclaimedJobs locks all currently unlocked jobs and returns all jobs locked by this process
func (o *orm) ClaimUnclaimedJobs(ctx context.Context) ([]models.JobSpecV2, error) {
	o.claimedJobsMu.Lock()
	defer o.claimedJobsMu.Unlock()

	claimedJobIDs := o.claimedJobIDs()

	var join string
	var args []interface{}
	if len(claimedJobIDs) > 0 {
		// NOTE: OFFSET 0 is a postgres trick that doesn't change the result,
		// but prevents the optimiser from trying to pull the where condition
		// up out of the subquery
		join = `
            INNER JOIN (
                SELECT not_claimed_by_us.id, pg_try_advisory_lock(?::integer, not_claimed_by_us.id) AS locked
                FROM (SELECT id FROM jobs WHERE id != ANY(?) OFFSET 0) not_claimed_by_us
            ) claimed_jobs ON jobs.id = claimed_jobs.id AND claimed_jobs.locked
        `
		args = []interface{}{o.advisoryLockClassID, pq.Array(claimedJobIDs)}
	} else {
		join = `
            INNER JOIN (
                SELECT not_claimed_by_us.id, pg_try_advisory_lock(?::integer, not_claimed_by_us.id) AS locked
                FROM jobs not_claimed_by_us
            ) claimed_jobs ON jobs.id = claimed_jobs.id AND claimed_jobs.locked
        `
		args = []interface{}{o.advisoryLockClassID}
	}

	var newlyClaimedJobs []models.JobSpecV2
	err := o.db.
		Joins(join, args...).
		Preload("OffchainreportingOracleSpec").
		Find(&newlyClaimedJobs).Error
	if err != nil {
		return nil, errors.Wrap(err, "ClaimUnclaimedJobs failed to load jobs")
	}

	o.claimedJobs = append(o.claimedJobs, newlyClaimedJobs...)

	return newlyClaimedJobs, errors.Wrap(err, "Job Spawner ORM could not load unclaimed job specs")
}

func (o *orm) claimedJobIDs() (ids []int32) {
	ids = []int32{}
	for _, job := range o.claimedJobs {
		ids = append(ids, job.ID)
	}
	return
}

func (o *orm) CreateJob(ctx context.Context, jobSpec *models.JobSpecV2, taskDAG pipeline.TaskDAG) error {
	if taskDAG.HasCycles() {
		return errors.New("task DAG has cycles, which are not permitted")
	}

	ctx, cancel := utils.CombinedContext(ctx, o.config.DatabaseMaximumTxDuration())
	defer cancel()

	return postgres.GormTransaction(ctx, o.db, func(tx *gorm.DB) error {
		pipelineSpecID, err := o.pipelineORM.CreateSpec(ctx, taskDAG)
		if err != nil {
			return errors.Wrap(err, "failed to create pipeline spec")
		}
		jobSpec.PipelineSpecID = pipelineSpecID

		err = tx.Create(jobSpec).Error
		return errors.Wrap(err, "failed to create job")
	})
}

// DeleteJob removes a job that is claimed by this orm
// TODO: Extend this in future so it can delete any job and other nodes handle
// it gracefully
// See: https://www.pivotaltracker.com/story/show/175287919
func (o *orm) DeleteJob(ctx context.Context, id int32) error {
	o.claimedJobsMu.Lock()
	defer o.claimedJobsMu.Unlock()

	for i, j := range o.claimedJobs {
		if j.ID == id {
			// Delete the current job from the claimedJobs list
			o.claimedJobs[i] = o.claimedJobs[len(o.claimedJobs)-1] // Copy last element to current position
			o.claimedJobs = o.claimedJobs[:len(o.claimedJobs)-1]   // Truncate slice.

			err := o.advisoryLocker.Unlock(ctx, o.advisoryLockClassID, id)
			if err != nil {
				return errors.Wrap(err, "DeleteJob failed to unlock job")
			}
		}
	}

	// FIXME: Why not simply call jobORM.pipleineORM.DeleteJob?
	err := o.db.Exec(`
            WITH deleted_jobs AS (
            	DELETE FROM jobs WHERE id = $1 RETURNING offchainreporting_oracle_spec_id, pipeline_spec_id
            ),
            deleted_oracle_specs AS (
				DELETE FROM offchainreporting_oracle_specs WHERE id IN (SELECT offchainreporting_oracle_spec_id FROM deleted_jobs)
			)
			DELETE FROM pipeline_specs WHERE id IN (SELECT pipeline_spec_id FROM deleted_jobs)
    	`, id).Error
	return errors.Wrap(err, "DeleteJob failed to delete job")
}
