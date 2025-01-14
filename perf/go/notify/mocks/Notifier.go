// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	alerts "go.skia.org/infra/perf/go/alerts"
	clustering2 "go.skia.org/infra/perf/go/clustering2"

	context "context"

	frame "go.skia.org/infra/perf/go/ui/frame"

	mock "github.com/stretchr/testify/mock"

	provider "go.skia.org/infra/perf/go/git/provider"

	testing "testing"
)

// Notifier is an autogenerated mock type for the Notifier type
type Notifier struct {
	mock.Mock
}

// ExampleSend provides a mock function with given fields: ctx, alert
func (_m *Notifier) ExampleSend(ctx context.Context, alert *alerts.Alert) error {
	ret := _m.Called(ctx, alert)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *alerts.Alert) error); ok {
		r0 = rf(ctx, alert)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RegressionFound provides a mock function with given fields: ctx, commit, previousCommit, alert, cl, _a5
func (_m *Notifier) RegressionFound(ctx context.Context, commit provider.Commit, previousCommit provider.Commit, alert *alerts.Alert, cl *clustering2.ClusterSummary, _a5 *frame.FrameResponse) (string, error) {
	ret := _m.Called(ctx, commit, previousCommit, alert, cl, _a5)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, provider.Commit, provider.Commit, *alerts.Alert, *clustering2.ClusterSummary, *frame.FrameResponse) string); ok {
		r0 = rf(ctx, commit, previousCommit, alert, cl, _a5)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, provider.Commit, provider.Commit, *alerts.Alert, *clustering2.ClusterSummary, *frame.FrameResponse) error); ok {
		r1 = rf(ctx, commit, previousCommit, alert, cl, _a5)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RegressionMissing provides a mock function with given fields: ctx, commit, previousCommit, alert, cl, _a5, threadingReference
func (_m *Notifier) RegressionMissing(ctx context.Context, commit provider.Commit, previousCommit provider.Commit, alert *alerts.Alert, cl *clustering2.ClusterSummary, _a5 *frame.FrameResponse, threadingReference string) error {
	ret := _m.Called(ctx, commit, previousCommit, alert, cl, _a5, threadingReference)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, provider.Commit, provider.Commit, *alerts.Alert, *clustering2.ClusterSummary, *frame.FrameResponse, string) error); ok {
		r0 = rf(ctx, commit, previousCommit, alert, cl, _a5, threadingReference)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewNotifier creates a new instance of Notifier. It also registers a cleanup function to assert the mocks expectations.
func NewNotifier(t testing.TB) *Notifier {
	mock := &Notifier{}

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
