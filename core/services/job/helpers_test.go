package job

import "github.com/smartcontractkit/chainlink/core/store/models"

func GetORMAdvisoryLockClassID(oi ORM) int32 {
	return oi.(*orm).advisoryLockClassID
}

func GetORMClaimedJobs(oi ORM) []models.JobSpecV2 {
	o := oi.(*orm)
	o.claimedJobsMu.Lock()
	defer o.claimedJobsMu.Unlock()
	copied := make([]models.JobSpecV2, len(o.claimedJobs))
	for i, j := range o.claimedJobs {
		copied[i] = j
	}
	return copied
}
