package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/akamensky/argparse"
	"github.com/matthiasng/service-wrapper/cli"
	"github.com/matthiasng/service-wrapper/service"
	"golang.org/x/sys/windows/svc"
)

func main() {
	parser := argparse.NewParser("service-wrapper", `Run a "-command" with "-arguments" as service`)

	serviceName := parser.String("n", "name", &argparse.Options{
		Required: true,
		Help:     "Servicename",
	})
	logDirectory := parser.String("l", "logdirectory", &argparse.Options{
		Required: true,
		Help:     "Log directory",
	})
	command := parser.String("c", "command", &argparse.Options{
		Required: true,
		Help:     "Command",
	})
	arguments := parser.List("a", "arg", &argparse.Options{
		Required: true,
		Help:     `Command arguments. Example: '... -a "-key" -a "value" -a "--key2" -a "value"'`,
	})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(2)
	}

	//
	isAnInteractiveSession, err := svc.IsAnInteractiveSession()
	if err != nil {
		log.Fatalf("failed to determine if we are running in an interactive session: %v", err)
	}

	//
	logger := service.Logger{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	if isAnInteractiveSession {
		err = os.MkdirAll(*logDirectory, os.ModePerm)
		if err != nil {
			log.Fatalf("Error - os.MkdirAll(%s, os.ModePerm): %v", *logDirectory, err)
		}

		logFileName := fmt.Sprintf("%s_%s.log", *serviceName, time.Now().Format("02-01-2006_15-04-05"))
		logFilePath := path.Join(*logDirectory, logFileName)
		file, err := os.Create(logFilePath)
		if err != nil {
			log.Fatalf("Error - os.Create(%s): %v", logFilePath, err)
		}
		defer func() { _ = file.Close() }()

		logger.Stdout = file
		logger.Stderr = file
	}

	//
	wrapper := service.Wrapper{
		Config: service.Configuration{
			IsAnInteractiveSession: isAnInteractiveSession,
			ServiceName:            *serviceName,
			Command:                *command,
			Arguments:              cli.BindArguments(*arguments),
			Logger:                 logger,
		},
	}
	err = wrapper.Run()
	if err != nil {
		log.Fatalf("wrapper.Run(): %v", err)
	}
}
