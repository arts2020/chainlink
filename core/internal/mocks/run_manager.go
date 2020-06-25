// Code generated by mockery v2.0.0. DO NOT EDIT.

package mocks

import (
	big "math/big"

	models "github.com/smartcontractkit/chainlink/core/store/models"
	mock "github.com/stretchr/testify/mock"
)

// RunManager is an autogenerated mock type for the RunManager type
type RunManager struct {
	mock.Mock
}

// Cancel provides a mock function with given fields: runID
func (_m *RunManager) Cancel(runID *models.ID) (*models.JobRun, error) {
	ret := _m.Called(runID)

	var r0 *models.JobRun
	if rf, ok := ret.Get(0).(func(*models.ID) *models.JobRun); ok {
		r0 = rf(runID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.JobRun)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*models.ID) error); ok {
		r1 = rf(runID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Create provides a mock function with given fields: jobSpecID, initiator, creationHeight, runRequest
func (_m *RunManager) Create(jobSpecID *models.ID, initiator *models.Initiator, creationHeight *big.Int, runRequest *models.RunRequest) (*models.JobRun, error) {
	ret := _m.Called(jobSpecID, initiator, creationHeight, runRequest)

	var r0 *models.JobRun
	if rf, ok := ret.Get(0).(func(*models.ID, *models.Initiator, *big.Int, *models.RunRequest) *models.JobRun); ok {
		r0 = rf(jobSpecID, initiator, creationHeight, runRequest)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.JobRun)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*models.ID, *models.Initiator, *big.Int, *models.RunRequest) error); ok {
		r1 = rf(jobSpecID, initiator, creationHeight, runRequest)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateErrored provides a mock function with given fields: jobSpecID, initiator, err
func (_m *RunManager) CreateErrored(jobSpecID *models.ID, initiator models.Initiator, err error) (*models.JobRun, error) {
	ret := _m.Called(jobSpecID, initiator, err)

	var r0 *models.JobRun
	if rf, ok := ret.Get(0).(func(*models.ID, models.Initiator, error) *models.JobRun); ok {
		r0 = rf(jobSpecID, initiator, err)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.JobRun)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*models.ID, models.Initiator, error) error); ok {
		r1 = rf(jobSpecID, initiator, err)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ResumeAllInProgress provides a mock function with given fields:
func (_m *RunManager) ResumeAllInProgress() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ResumeAllPendingConnection provides a mock function with given fields:
func (_m *RunManager) ResumeAllPendingConnection() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ResumeAllPendingNextBlock provides a mock function with given fields: currentBlockHeight
func (_m *RunManager) ResumeAllPendingNextBlock(currentBlockHeight *big.Int) error {
	ret := _m.Called(currentBlockHeight)

	var r0 error
	if rf, ok := ret.Get(0).(func(*big.Int) error); ok {
		r0 = rf(currentBlockHeight)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ResumePendingBridge provides a mock function with given fields: runID, input
func (_m *RunManager) ResumePendingBridge(runID *models.ID, input models.BridgeRunResult) error {
	ret := _m.Called(runID, input)

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.ID, models.BridgeRunResult) error); ok {
		r0 = rf(runID, input)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
