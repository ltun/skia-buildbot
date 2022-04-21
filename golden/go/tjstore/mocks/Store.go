// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	continuous_integration "go.skia.org/infra/golden/go/continuous_integration"

	testing "testing"

	time "time"

	tjstore "go.skia.org/infra/golden/go/tjstore"
)

// Store is an autogenerated mock type for the Store type
type Store struct {
	mock.Mock
}

// GetResults provides a mock function with given fields: ctx, psID, updatedAfter
func (_m *Store) GetResults(ctx context.Context, psID tjstore.CombinedPSID, updatedAfter time.Time) ([]tjstore.TryJobResult, error) {
	ret := _m.Called(ctx, psID, updatedAfter)

	var r0 []tjstore.TryJobResult
	if rf, ok := ret.Get(0).(func(context.Context, tjstore.CombinedPSID, time.Time) []tjstore.TryJobResult); ok {
		r0 = rf(ctx, psID, updatedAfter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]tjstore.TryJobResult)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, tjstore.CombinedPSID, time.Time) error); ok {
		r1 = rf(ctx, psID, updatedAfter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTryJob provides a mock function with given fields: ctx, id, cisName
func (_m *Store) GetTryJob(ctx context.Context, id string, cisName string) (continuous_integration.TryJob, error) {
	ret := _m.Called(ctx, id, cisName)

	var r0 continuous_integration.TryJob
	if rf, ok := ret.Get(0).(func(context.Context, string, string) continuous_integration.TryJob); ok {
		r0 = rf(ctx, id, cisName)
	} else {
		r0 = ret.Get(0).(continuous_integration.TryJob)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, id, cisName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTryJobs provides a mock function with given fields: ctx, psID
func (_m *Store) GetTryJobs(ctx context.Context, psID tjstore.CombinedPSID) ([]continuous_integration.TryJob, error) {
	ret := _m.Called(ctx, psID)

	var r0 []continuous_integration.TryJob
	if rf, ok := ret.Get(0).(func(context.Context, tjstore.CombinedPSID) []continuous_integration.TryJob); ok {
		r0 = rf(ctx, psID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]continuous_integration.TryJob)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, tjstore.CombinedPSID) error); ok {
		r1 = rf(ctx, psID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewStore creates a new instance of Store. It also registers a cleanup function to assert the mocks expectations.
func NewStore(t testing.TB) *Store {
	mock := &Store{}

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
