package main

import (
	"fmt"
	"log"
	"os"

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
	wrapper := service.Wrapper{
		Config: service.Configuration{
			ServiceName:            *serviceName,
			Command:                *command,
			Arguments:              cli.BindArguments(*arguments),
			LogDirectory:           *logDirectory,
			IsAnInteractiveSession: isAnInteractiveSession,
		},
	}
	err = wrapper.Run()
	if err != nil {
		log.Fatalf("wrapper.Run(): %v", err)
	}
}
