// Package targetconnect initiates and maintains a connection from
// a test machine to a switchboard pod. See https://go/skia-switchboard.
package targetconnect

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.skia.org/infra/go/testutils"
	"go.skia.org/infra/go/testutils/unittest"
	rpfMock "go.skia.org/infra/machine/go/machine/targetconnect/mocks"
	"go.skia.org/infra/machine/go/switchboard"
	"go.skia.org/infra/machine/go/switchboard/mocks"
)

const (
	hostname = "skia-rpi2-rack4-shelf1-002"
	username = "root"
)

var (
	errMyMockError = errors.New("my mock error")
	meetingPoint   = switchboard.MeetingPoint{
		PodName: "switch-pod-0",
		Port:    123,
	}
)

func TestSingleStep_FirstCallToReserveMeetingPointReturnsError_Returns(t *testing.T) {
	unittest.SmallTest(t)
	switchboardMock := &mocks.Switchboard{}
	switchboardMock.On("ReserveMeetingPoint", testutils.AnyContext, hostname, username).Return(meetingPoint, errMyMockError)
	rpf := &rpfMock.RevPortForward{}

	c := New(switchboardMock, rpf, hostname, username)
	c.singleStep(context.Background(), time.NewTicker(time.Microsecond), time.Microsecond)
	switchboardMock.AssertExpectations(t)
}

func TestSingleStep_KeepAliveMeetingPointGetsCalledMultipleTimes_Returns(t *testing.T) {
	unittest.SmallTest(t)
	ctx, cancel := context.WithCancel(context.Background())

	switchboardMock := &mocks.Switchboard{}
	switchboardMock.On("ReserveMeetingPoint", testutils.AnyContext, hostname, username).Return(meetingPoint, nil)
	keepAliveCount := 0
	switchboardMock.On("KeepAliveMeetingPoint", testutils.AnyContext, meetingPoint).Run(func(args mock.Arguments) {
		keepAliveCount++
		if keepAliveCount > 1 {
			cancel()
		}
	}).Times(2).Return(nil)
	switchboardMock.On("ClearMeetingPoint", testutils.AnyContext, meetingPoint).Return(nil)
	rpf := &rpfMock.RevPortForward{}
	rpf.On("Start", testutils.AnyContext, meetingPoint.PodName, meetingPoint.Port).Run(func(args mock.Arguments) {
		<-ctx.Done()
	}).Return(nil)

	c := New(switchboardMock, rpf, hostname, username)
	c.singleStep(ctx, time.NewTicker(time.Millisecond), time.Microsecond)
	switchboardMock.AssertExpectations(t)
}

func TestStart_ContextIsCancelled_ReturnsAndMeetingPointIsCleared(t *testing.T) {
	unittest.SmallTest(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	switchboardMock := &mocks.Switchboard{}
	switchboardMock.On("ReserveMeetingPoint", testutils.AnyContext, hostname, username).Return(meetingPoint, nil)
	var clearMeetingPointCalledWG sync.WaitGroup
	clearMeetingPointCalledWG.Add(1)
	switchboardMock.On("ClearMeetingPoint", testutils.AnyContext, meetingPoint).Run(func(mock.Arguments) {
		clearMeetingPointCalledWG.Done()
	}).Return(nil)
	rpf := &rpfMock.RevPortForward{}
	rpf.On("Start", testutils.AnyContext, meetingPoint.PodName, meetingPoint.Port).Run(func(args mock.Arguments) {
		<-ctx.Done()
	}).Return(nil)

	c := New(switchboardMock, rpf, hostname, username)
	err := c.Start(ctx)
	require.Error(t, err)
	clearMeetingPointCalledWG.Wait()
	require.Equal(t, int64(1), c.stepsCounter.Get())
	c.stepsCounter.Reset()
	switchboardMock.AssertExpectations(t)
}

func TestStart_FirstCallToRevPortForwardFails_CausesASecondCalltoConnectToPod(t *testing.T) {
	unittest.SmallTest(t)
	ctx, cancel := context.WithCancel(context.Background())
	switchboardMock := &mocks.Switchboard{}
	switchboardMock.On("ClearMeetingPoint", testutils.AnyContext, meetingPoint).Times(2).Return(nil)
	var reserveMeetingPointWG sync.WaitGroup
	reserveMeetingPointWG.Add(2)
	switchboardMock.On("ReserveMeetingPoint", testutils.AnyContext, hostname, username).Times(2).Run(func(args mock.Arguments) {
		reserveMeetingPointWG.Done()
	}).Return(meetingPoint, nil)
	rpf := &rpfMock.RevPortForward{}
	rpf.On("Start", testutils.AnyContext, meetingPoint.PodName, meetingPoint.Port).Return(errMyMockError).Times(1)
	rpf.On("Start", testutils.AnyContext, meetingPoint.PodName, meetingPoint.Port).Run(func(args mock.Arguments) {
		<-ctx.Done()
	}).Return(errMyMockError).Times(1)

	c := New(switchboardMock, rpf, hostname, username)
	// Call Start() in a Go routine since we need to cancel the Context after
	// Start() is called, and Start() doesn't return.
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := c.Start(ctx)
		require.Error(t, err)
	}()

	// Wait until ReserveMeetingPoint has been called twice.
	reserveMeetingPointWG.Wait()
	cancel()

	// Wait for Start() to return.
	wg.Wait()
	require.Equal(t, int64(2), c.stepsCounter.Get())
	c.stepsCounter.Reset()
	switchboardMock.AssertExpectations(t)
}