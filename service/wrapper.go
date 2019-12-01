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

// Configuration for Run
type Configuration struct {
	ServiceName            string
	IsAnInteractiveSession bool
	Command                string
	Arguments              []string
	LogDirectory           string
}

// Wrapper handles command execution
type Wrapper struct {
	Config   Configuration
	eventLog debug.Log
	cmd      *exec.Cmd
}

// Run runs command
func (s *Wrapper) Run() error {
	var err error
	if s.Config.IsAnInteractiveSession {
		s.eventLog = debug.New(s.Config.ServiceName)
	} else {
		s.eventLog, err = eventlog.Open(s.Config.ServiceName) // #todo replace eventlog with custom logfile
		if err != nil {
			return err
		}
	}
	defer func() { _ = s.eventLog.Close() }()

	s.cmd = exec.Command(s.Config.Command, s.Config.Arguments...)

	var runFunc func(string, svc.Handler) error
	if s.Config.IsAnInteractiveSession {
		s.cmd.Stdout = os.Stdout
		s.cmd.Stderr = os.Stderr

		runFunc = debug.Run
	} else {
		err = os.MkdirAll(s.Config.LogDirectory, os.ModePerm)
		if err != nil {
			_ = s.eventLog.Error(1, fmt.Sprintf("Error - os.MkdirAll(%s, os.ModePerm): %v", s.Config.LogDirectory, err))
			return err
		}

		logFileName := fmt.Sprintf("%s_%s.log", s.Config.ServiceName, time.Now().Format("02-01-2006_15-04-05"))
		logFilePath := path.Join(s.Config.LogDirectory, logFileName)
		file, err := os.Create(logFilePath)
		if err != nil {
			_ = s.eventLog.Error(1, fmt.Sprintf("Error - os.Create(%s): %v", logFilePath, err))
			return err
		}
		defer func() { _ = file.Close() }()

		s.cmd.Stdout = file
		s.cmd.Stderr = file

		runFunc = svc.Run
	}

	winSvc := winService{
		wrapper: s,
	}

	_ = s.eventLog.Error(1, "service starting")
	err = runFunc(s.Config.ServiceName, &winSvc)
	if err != nil {
		_ = s.eventLog.Error(1, fmt.Sprintf("service failed: %v", err))
		return err
	}

	_ = s.eventLog.Info(1, "service stopped")
	return nil
}
