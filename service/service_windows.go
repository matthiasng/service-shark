package service

import (
	"os"
	"syscall"

	wsvc "golang.org/x/sys/windows/svc"
)

type windowsService struct {
	name          string
	program       Program
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

	if err := ws.program.Start(ws); err != nil {
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
				err := ws.program.Stop()
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
