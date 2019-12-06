// +build windows

package service

import (
	"os"
	"path/filepath"

	wsvc "golang.org/x/sys/windows/svc"
	wsvcDebug "golang.org/x/sys/windows/svc/debug"
)

// Create variables for svc and signal functions so we can mock them in tests
var svcIsAnInteractiveSession = wsvc.IsAnInteractiveSession

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
func Run(service Service, name string) error {
	var err error

	interactive, err := svcIsAnInteractiveSession()
	if err != nil {
		return err
	}

	ws := &windowsService{
		name:          name,
		svc:           service,
		isInteractive: interactive,
		exitChan:      make(chan error),
	}

	if ws.IsWindowsService() {
		// the working directory for a Windows Service is C:\Windows\System32
		// this is almost certainly not what the user wants.
		dir := filepath.Dir(os.Args[0])
		if err = os.Chdir(dir); err != nil {
			return err
		}
	}

	return ws.run()
}

type windowsService struct {
	name          string
	svc           Service
	isInteractive bool
	exitChan      chan error
	err           error
}

func (ws *windowsService) IsWindowsService() bool {
	return !ws.isInteractive
}

func (ws *windowsService) ExitService(err error) {
	ws.exitChan <- err
	close(ws.exitChan)
}

func (ws *windowsService) Name() string {
	return ws.name
}

func (ws *windowsService) setError(err error) {
	ws.err = err
}

func (ws *windowsService) getError() error {
	return ws.err
}

func (ws *windowsService) run() error {
	ws.setError(nil)

	run := wsvc.Run
	if !ws.IsWindowsService() {
		run = wsvcDebug.Run
	}

	runErr := run(ws.name, ws)
	startStopErr := ws.getError()
	if startStopErr != nil {
		return startStopErr
	}
	if runErr != nil {
		return runErr
	}
	return nil
}

// Execute is invoked by Windows
func (ws *windowsService) Execute(args []string, r <-chan wsvc.ChangeRequest, changes chan<- wsvc.Status) (bool, uint32) {
	changes <- wsvc.Status{State: wsvc.StartPending}

	if err := ws.svc.Start(ws); err != nil {
		ws.setError(err)
		return true, 1
	}

	changes <- wsvc.Status{State: wsvc.Running, Accepts: wsvc.AcceptStop | wsvc.AcceptShutdown}
loop:
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case wsvc.Interrogate:
				changes <- c.CurrentStatus
			case wsvc.Stop, wsvc.Shutdown:
				changes <- wsvc.Status{State: wsvc.StopPending}
				err := ws.svc.Stop()
				if err != nil {
					ws.setError(err)
					return true, 2
				}
				break loop
			default:
				continue loop
			}
		case err := <-ws.exitChan:
			ws.setError(err)
			return true, 3
		}
	}

	return false, 0
}
