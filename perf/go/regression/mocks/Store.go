// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	clustering2 "go.skia.org/infra/perf/go/clustering2"

	frame "go.skia.org/infra/perf/go/ui/frame"

	mock "github.com/stretchr/testify/mock"

	regression "go.skia.org/infra/perf/go/regression"

	testing "testing"

	types "go.skia.org/infra/perf/go/types"
)

// Store is an autogenerated mock type for the Store type
type Store struct {
	mock.Mock
}

// Range provides a mock function with given fields: ctx, begin, end
func (_m *Store) Range(ctx context.Context, begin types.CommitNumber, end types.CommitNumber) (map[types.CommitNumber]*regression.AllRegressionsForCommit, error) {
	ret := _m.Called(ctx, begin, end)

	var r0 map[types.CommitNumber]*regression.AllRegressionsForCommit
	if rf, ok := ret.Get(0).(func(context.Context, types.CommitNumber, types.CommitNumber) map[types.CommitNumber]*regression.AllRegressionsForCommit); ok {
		r0 = rf(ctx, begin, end)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[types.CommitNumber]*regression.AllRegressionsForCommit)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, types.CommitNumber, types.CommitNumber) error); ok {
		r1 = rf(ctx, begin, end)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetHigh provides a mock function with given fields: ctx, commitNumber, alertID, df, high
func (_m *Store) SetHigh(ctx context.Context, commitNumber types.CommitNumber, alertID string, df *frame.FrameResponse, high *clustering2.ClusterSummary) (bool, error) {
	ret := _m.Called(ctx, commitNumber, alertID, df, high)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, types.CommitNumber, string, *frame.FrameResponse, *clustering2.ClusterSummary) bool); ok {
		r0 = rf(ctx, commitNumber, alertID, df, high)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, types.CommitNumber, string, *frame.FrameResponse, *clustering2.ClusterSummary) error); ok {
		r1 = rf(ctx, commitNumber, alertID, df, high)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetLow provides a mock function with given fields: ctx, commitNumber, alertID, df, low
func (_m *Store) SetLow(ctx context.Context, commitNumber types.CommitNumber, alertID string, df *frame.FrameResponse, low *clustering2.ClusterSummary) (bool, error) {
	ret := _m.Called(ctx, commitNumber, alertID, df, low)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, types.CommitNumber, string, *frame.FrameResponse, *clustering2.ClusterSummary) bool); ok {
		r0 = rf(ctx, commitNumber, alertID, df, low)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, types.CommitNumber, string, *frame.FrameResponse, *clustering2.ClusterSummary) error); ok {
		r1 = rf(ctx, commitNumber, alertID, df, low)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TriageHigh provides a mock function with given fields: ctx, commitNumber, alertID, tr
func (_m *Store) TriageHigh(ctx context.Context, commitNumber types.CommitNumber, alertID string, tr regression.TriageStatus) error {
	ret := _m.Called(ctx, commitNumber, alertID, tr)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, types.CommitNumber, string, regression.TriageStatus) error); ok {
		r0 = rf(ctx, commitNumber, alertID, tr)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TriageLow provides a mock function with given fields: ctx, commitNumber, alertID, tr
func (_m *Store) TriageLow(ctx context.Context, commitNumber types.CommitNumber, alertID string, tr regression.TriageStatus) error {
	ret := _m.Called(ctx, commitNumber, alertID, tr)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, types.CommitNumber, string, regression.TriageStatus) error); ok {
		r0 = rf(ctx, commitNumber, alertID, tr)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Write provides a mock function with given fields: ctx, regressions
func (_m *Store) Write(ctx context.Context, regressions map[types.CommitNumber]*regression.AllRegressionsForCommit) error {
	ret := _m.Called(ctx, regressions)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, map[types.CommitNumber]*regression.AllRegressionsForCommit) error); ok {
		r0 = rf(ctx, regressions)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewStore creates a new instance of Store. It also registers a cleanup function to assert the mocks expectations.
func NewStore(t testing.TB) *Store {
	mock := &Store{}

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
