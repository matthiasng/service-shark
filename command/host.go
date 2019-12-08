package command

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/matthiasng/service-shark/cli"
	"github.com/matthiasng/service-shark/service"
)

type Host struct {
	Arguments     cli.Arguments
	cmd           *exec.Cmd
	logFile       *os.File
	quitSignal    chan struct{}
	quitCompleted chan struct{}
}

func (h *Host) init(env service.Environment) error {
	h.cmd = exec.Command(h.Arguments.Command, h.Arguments.CommandArguments...)

	if env.IsWindowsService() {
		err := os.MkdirAll(h.Arguments.LogDirectory, os.ModePerm)
		if err != nil {
			return err
		}

		logFileName := fmt.Sprintf("%s_%s.log", h.Name(), time.Now().Format("02-01-2006_15-04-05"))
		logFilePath := path.Join(h.Arguments.LogDirectory, logFileName)

		logFile, err := os.Create(logFilePath)
		if err != nil {
			return err
		}

		h.cmd.Stdout = logFile
		h.cmd.Stderr = logFile
	} else {
		h.cmd.Stdout = os.Stdout
		h.cmd.Stderr = os.Stderr
	}

	return h.cmd.Start()
}

func (h *Host) Start(env service.Environment) error {
	err := h.init(env)
	if err != nil {
		return err
	}

	h.quitSignal = make(chan struct{})
	h.quitCompleted = make(chan struct{})

	go func() {
		<-h.quitSignal
		_ = h.cmd.Process.Kill()
		close(h.quitCompleted)
	}()

	go func() {
		err = h.cmd.Wait()
		env.ExitService(err) // service command stopped -> service error
	}()

	return nil
}

func (h *Host) Stop() error {
	close(h.quitSignal)
	<-h.quitCompleted

	_ = h.logFile.Close()

	return nil
}

func (h *Host) Name() string {
	return h.Arguments.Name
}
