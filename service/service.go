package service

import (
	"fmt"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
	"os"
	"os/exec"
)

var (
	//command = "C:/Program Files/PowerShell/7-preview/pwsh.exe"
	command = "C:/Program Files/PowerShell/7-preview/preview/pwsh-preview.cmd"
	arguments = []string {
		"P:/_dev/projects/service-wrapper/test.ps1",
	}

	logFilePath = "P:/_dev/projects/service-wrapper/test.log"
)

var (
	eventLog debug.Log
	cmd *exec.Cmd
)



func Run(name string, isAnInteractiveSession bool) {
	var err error
	if isAnInteractiveSession {
		eventLog = debug.New(name)
	} else {
		eventLog, err = eventlog.Open(name)
		if err != nil {
			return
		}
	}
	defer func(){ _ = eventLog.Close() }()

	cmd = exec.Command(command, arguments...)

	var runFunc func(string, svc.Handler) error
	if isAnInteractiveSession {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		runFunc = debug.Run
	} else {
		file, err := os.Create(logFilePath)
		if err != nil {
			_ = eventLog.Error(1, fmt.Sprintf("Error - os.Create(%s): %v", logFilePath, err))
			return
		}
		defer func(){ _ = file.Close() }()

		cmd.Stdout = file
		cmd.Stderr = file

		runFunc = svc.Run
	}

	_ = eventLog.Error(1, fmt.Sprintf("%s service starting", name))
	err = runFunc(name, &serviceWrapper{})
	if err != nil {
		_ = eventLog.Error(1, fmt.Sprintf("%s service failed: %v", name, err))
		return
	}
	_ = eventLog.Info(1, fmt.Sprintf("%s service stopped", name))
}

type serviceWrapper struct{}

func (m *serviceWrapper) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (svcSpecificEC bool, exitCode  uint32) {
	changes <- svc.Status{State: svc.StartPending}

	//
	err := cmd.Start()
	if err != nil {
		_ = eventLog.Error(1, fmt.Sprintf("Error - cmd.Start(): %v", err))
		return true, 1
	}

	go func() {
		_ = eventLog.Error(1, "before wait")
		err = cmd.Wait()
		_ = eventLog.Error(1, "after wait")

		if err != nil {
			_ = eventLog.Error(1, fmt.Sprintf("Error - cmd.Wait(): %v", err))
			// #todo crash service!
		}

		os.Exit(1)
	} ()

	changes <- svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown}

loop:
	for {
		c := <-r
		switch c.Cmd {
		case svc.Interrogate:
			changes <- c.CurrentStatus
		case svc.Stop, svc.Shutdown:
			changes <- svc.Status{State: svc.StopPending}

			_ = eventLog.Info(1, "Kill process")
			_ = cmd.Process.Kill()	// #todo kill child process ? test with "C:/Program Files/PowerShell/7-preview/preview/pwsh-preview.cmd"
			_ = eventLog.Info(1, "dead!")

			changes <- svc.Status{State: svc.Stopped}
			break loop
		default:
			continue loop
		}
	}

	return false, 0
}
