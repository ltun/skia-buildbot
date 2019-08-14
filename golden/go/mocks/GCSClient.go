// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import io "io"
import mock "github.com/stretchr/testify/mock"
import storage "go.skia.org/infra/golden/go/storage"
import types "go.skia.org/infra/golden/go/types"

// GCSClient is an autogenerated mock type for the GCSClient type
type GCSClient struct {
	mock.Mock
}

// LoadKnownDigests provides a mock function with given fields: w
func (_m *GCSClient) LoadKnownDigests(w io.Writer) error {
	ret := _m.Called(w)

	var r0 error
	if rf, ok := ret.Get(0).(func(io.Writer) error); ok {
		r0 = rf(w)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Options provides a mock function with given fields:
func (_m *GCSClient) Options() storage.GCSClientOptions {
	ret := _m.Called()

	var r0 storage.GCSClientOptions
	if rf, ok := ret.Get(0).(func() storage.GCSClientOptions); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(storage.GCSClientOptions)
	}

	return r0
}

// RemoveForTestingOnly provides a mock function with given fields: targetPath
func (_m *GCSClient) RemoveForTestingOnly(targetPath string) error {
	ret := _m.Called(targetPath)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(targetPath)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WriteKnownDigests provides a mock function with given fields: digests
func (_m *GCSClient) WriteKnownDigests(digests types.DigestSlice) error {
	ret := _m.Called(digests)

	var r0 error
	if rf, ok := ret.Get(0).(func(types.DigestSlice) error); ok {
		r0 = rf(digests)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
