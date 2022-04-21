// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"
	testing "testing"

	mock "github.com/stretchr/testify/mock"
)

// ProjectMirror is an autogenerated mock type for the ProjectMirror type
type ProjectMirror struct {
	mock.Mock
}

// AddToDownloadedPackageTarballs provides a mock function with given fields: packageTarballName
func (_m *ProjectMirror) AddToDownloadedPackageTarballs(packageTarballName string) {
	_m.Called(packageTarballName)
}

// GetProjectName provides a mock function with given fields:
func (_m *ProjectMirror) GetProjectName() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// IsPackageTarballDownloaded provides a mock function with given fields: packageTarballName
func (_m *ProjectMirror) IsPackageTarballDownloaded(packageTarballName string) bool {
	ret := _m.Called(packageTarballName)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(packageTarballName)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// StartMirror provides a mock function with given fields: ctx, port
func (_m *ProjectMirror) StartMirror(ctx context.Context, port int) error {
	ret := _m.Called(ctx, port)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int) error); ok {
		r0 = rf(ctx, port)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewProjectMirror creates a new instance of ProjectMirror. It also registers a cleanup function to assert the mocks expectations.
func NewProjectMirror(t testing.TB) *ProjectMirror {
	mock := &ProjectMirror{}

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
