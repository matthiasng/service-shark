// +build windows

package service

import (
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	wsvc "golang.org/x/sys/windows/svc"
)

// Create variables for svc and signal functions so we can mock them in tests
var svcIsAnInteractiveSession = wsvc.IsAnInteractiveSession
var svcRun = wsvc.Run
var signalNotify = signal.Notify

// FixWorkingDirectory changes the working directory to the exeutable directory.
// The working directory for a Windows Service is C:\Windows\System32 ...
func FixWorkingDirectory() error {
	dir := filepath.Dir(os.Args[0])
	return os.Chdir(dir)
}

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
func Run(service Service) error {
	var err error

	interactive, err := svcIsAnInteractiveSession()
	if err != nil {
		return err
	}

	ws := &windowsService{
		name:          service.Name(),
		svc:           service,
		isInteractive: interactive,
		signalErrChan: make(chan error, 1),
	}

	return ws.run()
}

type windowsService struct {
	name          string
	svc           Service
	isInteractive bool
	signalErrChan chan error
	serviceErr    error
}

func (ws *windowsService) IsWindowsService() bool {
	return !ws.isInteractive
}

func (ws *windowsService) ExitService(err error) {
	ws.signalErrChan <- err
}

func (ws *windowsService) run() error {
	var runErr error
	if ws.IsWindowsService() {
		runErr = svcRun(ws.name, ws)
	} else {
		runErr = ws.runInteractive(ws.name, ws)
	}

	if ws.serviceErr != nil {
		return ws.serviceErr
	}

	return runErr
}

func (ws *windowsService) runInteractive(name string, handler wsvc.Handler) error {
	cmds := make(chan wsvc.ChangeRequest)
	changes := make(chan wsvc.Status)

	sig := make(chan os.Signal)
	signalNotify(sig, os.Interrupt)

	go func() {
		status := wsvc.Status{State: wsvc.Stopped}
		for {
			select {
			case <-sig:
				cmds <- wsvc.ChangeRequest{Cmd: wsvc.Stop, CurrentStatus: status}
			case <-changes:
			}
		}
	}()

	_, runErrNo := handler.Execute([]string{name}, cmds, changes)
	if runErrNo != 0 {
		return syscall.Errno(runErrNo)
	}

	return nil
}

// Execute is invoked by Windows
func (ws *windowsService) Execute(args []string, r <-chan wsvc.ChangeRequest, changes chan<- wsvc.Status) (bool, uint32) {
	changes <- wsvc.Status{State: wsvc.StartPending}

	if err := ws.svc.Start(ws); err != nil {
		ws.serviceErr = err
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
					ws.serviceErr = err
					return true, 2
				}
				break loop
			default:
				continue loop
			}
		case err := <-ws.signalErrChan:
			ws.serviceErr = err
			return true, 3
		}
	}

	return false, 0
}
