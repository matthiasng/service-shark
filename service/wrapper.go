package service

import (
	"fmt"
	"io"
	"os/exec"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
)

// Configuration for Run
type Configuration struct {
	IsAnInteractiveSession bool
	ServiceName            string
	Command                string
	Arguments              []string
	Logger                 Logger
}

// Wrapper handles command execution
type Wrapper struct {
	Config   Configuration
	eventLog debug.Log
	cmd      *exec.Cmd
}

type Logger struct {
	Stdout io.Writer
	Stderr io.Writer
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

	fmt.Println(s.Config.Command, s.Config.Arguments)

	s.cmd = exec.Command(s.Config.Command, s.Config.Arguments...)
	s.cmd.Stdout = s.Config.Logger.Stdout
	s.cmd.Stderr = s.Config.Logger.Stderr

	winSvc := winService{
		wrapper: s,
	}

	_ = s.eventLog.Error(1, "service starting")

	if s.Config.IsAnInteractiveSession {
		err = debug.Run(s.Config.ServiceName, &winSvc)
	} else {
		err = svc.Run(s.Config.ServiceName, &winSvc)
	}

	if err != nil {
		_ = s.eventLog.Error(1, fmt.Sprintf("service failed: %v", err))
		return err
	}

	_ = s.eventLog.Info(1, "service stopped")
	return nil
}
