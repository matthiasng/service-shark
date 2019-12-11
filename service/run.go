package service

import (
	"os/signal"

	wsvc "golang.org/x/sys/windows/svc"
)

// Create variables for svc and signal functions so we can mock them in tests
var svcIsAnInteractiveSession = wsvc.IsAnInteractiveSession
var svcRun = wsvc.Run
var signalNotify = signal.Notify

// Run runs an implementation of the Service interface.
//
// Run will block until the Windows Service is stopped or Ctrl+C is pressed if
// running from the console.
//
// Stopping the Windows Service and Ctrl+C will call the Service's Stop method to
// initiate a graceful shutdown.
//
// Note that WM_CLOSE is not handled (end task) and the Service's Stop method will
// not be called.
func Run(program Program) error {
	var err error

	interactive, err := svcIsAnInteractiveSession()
	if err != nil {
		return err
	}

	ws := &windowsService{
		name:          program.Name(),
		program:       program,
		isInteractive: interactive,
		signalErrChan: make(chan error, 1),
	}

	return ws.run()
}
