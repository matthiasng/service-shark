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

// Host implements the service.Program interface to run any command as service
type Host struct {
	Arguments cli.Arguments
	cmd       *exec.Cmd
	logFile   *os.File
}

func (h *Host) init(env service.Environment) error {
	h.cmd = exec.Command(h.Arguments.Command, h.Arguments.CommandArguments...)

	if env.IsWindowsService() {
		err := os.MkdirAll(h.Arguments.LogDirectory, os.ModePerm)
		if err != nil {
			return err
		}

		logFileName := fmt.Sprintf("%s_%s.log", h.Name(), time.Now().Format("2006-01-02-_15-04-05"))
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

// Start prepares logging and executes the command
func (h *Host) Start(env service.Environment) error {
	err := h.init(env)
	if err != nil {
		return err
	}

	go func() {
		err = h.cmd.Wait()
		h.cleanup()
		env.ExitService(err) // service command stopped -> service error
	}()

	return nil
}

// Stop kills the command and closes the log file
func (h *Host) Stop() error {
	_ = h.cmd.Process.Kill()
	h.cleanup()

	return nil
}

func (h *Host) cleanup() {
	_ = h.logFile.Close()
}

func (h *Host) Name() string {
	return h.Arguments.Name
}
