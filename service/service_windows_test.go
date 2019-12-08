// +build windows

package service

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/sys/windows/svc"
)

type mockProgram struct {
	startCalled int
	stopCalled  int
	name        string
	start       func(Environment) error
	stop        func() error
}

func (p *mockProgram) Start(env Environment) error {
	p.startCalled++
	return p.start(env)
}

func (p *mockProgram) Stop() error {
	p.stopCalled++
	return p.stop()
}

func (p *mockProgram) Name() string {
	return p.name
}

func makeProgram() *mockProgram {
	return &mockProgram{
		start: func(Environment) error {
			return nil
		},
		stop: func() error {
			return nil
		},
	}
}

func mock_svcIsAnInteractiveSession(fnc func() (bool, error)) func() {
	o := svcIsAnInteractiveSession
	svcIsAnInteractiveSession = fnc
	return func() {
		svcIsAnInteractiveSession = o
	}
}

func mock_svcRun(fnc func(string, svc.Handler) error) func() {
	o := svcRun
	svcRun = fnc
	return func() {
		svcRun = o
	}
}

func mock_signalNotify(fnc func(c chan<- os.Signal, sig ...os.Signal)) func() {
	o := signalNotify
	signalNotify = fnc
	return func() {
		signalNotify = o
	}
}

func Test_IsAnInteractiveSession_Error(t *testing.T) {
	require := require.New(t)
	prg := makeProgram()

	testErr := errors.New("test error")

	defer mock_svcIsAnInteractiveSession(func() (bool, error) {
		return false, testErr
	})()

	err := Run(prg)
	require.Equal(err, testErr)
	require.Zero(prg.startCalled)
	require.Zero(prg.stopCalled)
}

func Test_svcRun_Call(t *testing.T) {
	require := require.New(t)
	prg := makeProgram()
	prg.name = "test"

	defer mock_svcIsAnInteractiveSession(func() (bool, error) {
		return false, nil
	})()

	svcName := ""
	svcRunCalled := false
	defer mock_svcRun(func(n string, h svc.Handler) error {
		svcName = n
		svcRunCalled = true
		return nil
	})()

	err := Run(prg)

	require.NoError(err)
	require.True(svcRunCalled)
	require.Equal(svcName, "test")
	require.Zero(prg.startCalled)
	require.Zero(prg.stopCalled)
}

func Test_Start_Error(t *testing.T) {
	require := require.New(t)
	prg := makeProgram()

	startErr := errors.New("start error")
	prg.start = func(Environment) error {
		return startErr
	}

	err := Run(prg)

	require.Equal(err, startErr)
}

func Test_Stop(t *testing.T) {
	require := require.New(t)
	prg := makeProgram()

	defer mock_signalNotify(func(c chan<- os.Signal, sig ...os.Signal) {
		go func() {
			time.Sleep(2 * time.Second)
			c <- os.Interrupt
		}()
	})()

	err := Run(prg)

	require.NoError(err)
	require.Equal(prg.startCalled, 1)
	require.Equal(prg.stopCalled, 1)
}

func Test_Stop_Error(t *testing.T) {
	require := require.New(t)
	prg := makeProgram()

	stopErr := errors.New("stop error")
	prg.stop = func() error {
		return stopErr
	}

	defer mock_signalNotify(func(c chan<- os.Signal, sig ...os.Signal) {
		go func() {
			time.Sleep(2 * time.Second)
			c <- os.Interrupt
		}()
	})()

	err := Run(prg)

	require.Equal(err, stopErr)
	require.Equal(prg.startCalled, 1)
	require.Equal(prg.stopCalled, 1)
}

func Test_ExitService(t *testing.T) {
	require := require.New(t)
	prg := makeProgram()

	exitErr := errors.New("exit error")
	prg.start = func(e Environment) error {
		go func() {
			time.Sleep(2 * time.Second)
			e.ExitService(exitErr)
		}()
		return nil
	}

	err := Run(prg)

	require.Equal(err, exitErr)
	require.Equal(prg.startCalled, 1)
	require.Equal(prg.stopCalled, 0)
}
