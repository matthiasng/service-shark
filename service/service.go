package service

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
)

var (
	eventLog debug.Log
	cmd      *exec.Cmd
)

type Configuration struct {
	ServiceName  string
	Command      string
	Arguments    []string
	LogDirectory string // #todo default
}

func Run(config Configuration, isAnInteractiveSession bool) {
	var err error
	if isAnInteractiveSession {
		eventLog = debug.New(config.ServiceName)
	} else {
		eventLog, err = eventlog.Open(config.ServiceName) // #todo replace eventlog with custom logfile
		if err != nil {
			return
		}
	}
	defer func() { _ = eventLog.Close() }()

	cmd = exec.Command(config.Command, config.Arguments...)

	var runFunc func(string, svc.Handler) error
	if isAnInteractiveSession {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		runFunc = debug.Run
	} else {
		err = os.MkdirAll(config.LogDirectory, os.ModePerm)
		if err != nil {
			_ = eventLog.Error(1, fmt.Sprintf("Error - os.MkdirAll(%s, os.ModePerm): %v", config.LogDirectory, err))
			return
		}

		logFileName := time.Now().Format("02-01-2006_15-04-05") + ".log"
		logFilePath := path.Join(config.LogDirectory, logFileName)
		file, err := os.Create(logFilePath)
		if err != nil {
			_ = eventLog.Error(1, fmt.Sprintf("Error - os.Create(%s): %v", logFilePath, err))
			return
		}
		defer func() { _ = file.Close() }()

		cmd.Stdout = file
		cmd.Stderr = file

		runFunc = svc.Run
	}

	_ = eventLog.Error(1, "service starting")
	err = runFunc(config.ServiceName, &serviceWrapper{})
	if err != nil {
		_ = eventLog.Error(1, fmt.Sprintf("service failed: %v", err))
		return
	}
	_ = eventLog.Info(1, "service stopped")
}

type serviceWrapper struct{}

func (m *serviceWrapper) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (svcSpecificEC bool, exitCode uint32) {
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
		}

		os.Exit(1)
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

			_ = eventLog.Info(1, "Kill process")
			_ = cmd.Process.Kill() // #todo kill child process. Test Command -> "C:/Program Files/PowerShell/7-preview/preview/pwsh-preview.cmd"
			_ = eventLog.Info(1, "dead!")

			changes <- svc.Status{State: svc.Stopped}
			break loop
		default:
			continue loop
		}
	}

	return false, 0
}
