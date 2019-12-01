package service

import (
	"fmt"
	"os"

	"golang.org/x/sys/windows/svc"
)

type winService struct {
	wrapper *Wrapper
}

func (w *winService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (svcSpecificEC bool, exitCode uint32) {
	changes <- svc.Status{State: svc.StartPending}

	err := w.wrapper.cmd.Start()
	if err != nil {
		_ = w.wrapper.eventLog.Error(1, fmt.Sprintf("Error - cmd.Start(): %v", err))
		return true, 1
	}

	go func() {
		err = w.wrapper.cmd.Wait()
		if err != nil {
			_ = w.wrapper.eventLog.Error(1, fmt.Sprintf("Error - cmd.Wait(): %v", err))
		}
		os.Exit(1) // error because our service command stopped
	}()

	changes <- svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown}

loop:
	for {
		c := <-r
		switch c.Cmd {
		case svc.Interrogate:
			changes <- c.CurrentStatus
		case svc.Stop, svc.Shutdown:
			changes <- svc.Status{State: svc.StopPending}

			err := w.wrapper.cmd.Process.Kill() // #todo kill child process. Test Command -> "C:/Program Files/PowerShell/7-preview/preview/pwsh-preview.cmd"
			if err != nil {
				_ = w.wrapper.eventLog.Error(1, fmt.Sprintf("Error - cmd.Process.Kill(): %v", err))
			}

			changes <- svc.Status{State: svc.Stopped}
			break loop
		default:
			continue loop
		}
	}

	return false, 0
}
