package command

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/matthiasng/service-shark/service"
)

type Config struct {
	Name      string
	Arguments []string
}

type Host struct {
	CmdConfig     Config
	LogDirecotry  string
	cmd           *exec.Cmd
	logFile       *os.File
	quitSignal    chan struct{}
	quitCompleted chan struct{}
}

func (h *Host) Init(env service.Environment) error {
	h.cmd = exec.Command(h.CmdConfig.Name, h.CmdConfig.Arguments...)

	if env.IsWindowsService() {
		err := os.MkdirAll(h.LogDirecotry, os.ModePerm)
		if err != nil {
			return err
		}

		name := "" // #todo

		logFileName := fmt.Sprintf("%s_%s.log", name, time.Now().Format("02-01-2006_15-04-05"))
		logFilePath := path.Join(h.LogDirecotry, logFileName)

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

func (h *Host) Start() error {
	h.quitSignal = make(chan struct{})
	h.quitCompleted = make(chan struct{})

	go func() {
		<-h.quitSignal
		_ = h.cmd.Process.Kill() // #todo kill child process. Test Command -> "C:/Program Files/PowerShell/7-preview/preview/pwsh-preview.cmd"
		close(h.quitCompleted)
	}()

	go func() {
		_ = h.cmd.Wait()
		os.Exit(1) // service command stopped -> stop service
	}()

	return nil
}

func (h *Host) Stop() error {
	close(h.quitSignal)
	<-h.quitCompleted

	_ = h.logFile.Close()

	return nil
}
