package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/bioapfelsaft/service-wrapper/service"
	"golang.org/x/sys/windows/svc"
)

// #todo config

func parseCommandLine(commandArgs string) ([]string, error) {
	args := []string{}

	insideQuotes := false
	quoteStartPos := 0

	currentArg := ""

	for i := 0; i < len(commandArgs); i++ {
		c := commandArgs[i]

		if c == '\'' {
			quoteStartPos = i
			if insideQuotes && (i+1) < len(commandArgs) && commandArgs[i+1] != ' ' {
				return []string{}, fmt.Errorf("end quote must be followed by a space - position %d", i)
			}

			insideQuotes = !insideQuotes
			continue
		}

		if c == ' ' && !insideQuotes {
			if len(currentArg) > 0 {
				args = append(args, currentArg)
				currentArg = ""
			}

			continue
		}

		currentArg += string(c)
	}

	if insideQuotes {
		return []string{}, fmt.Errorf("Unclosed quote in command line - start position %d", quoteStartPos)
	}

	args = append(args, currentArg)
	return args, nil
}

func main() {
	serviceName := "WinServiceWrapper" // #todo ?

	command := flag.String("command", "", "Command")
	logDirectory := flag.String("logdir", "", "Log directory")
	arguments := flag.String("arguments", "", "Arguments")

	flag.Parse()

	args, err := parseCommandLine(*arguments)
	if err != nil {
		log.Fatalf("invalid arguments: %v", err)
	}

	fmt.Println(*arguments)
	fmt.Println("-----------------")
	for _, e := range args {
		fmt.Println(e)
	}

	//
	isAnInteractiveSession, err := svc.IsAnInteractiveSession()
	if err != nil {
		log.Fatalf("failed to determine if we are running in an interactive session: %v", err)
	}

	//
	service.Run(service.Configuration{
		ServiceName:  serviceName,
		Command:      *command,
		Arguments:    args,
		LogDirectory: *logDirectory,
	}, isAnInteractiveSession)

	// service.Run(service.Configuration{
	// 	ServiceName: *serviceName,
	// 	Command:     "C:/Program Files/PowerShell/7-preview/pwsh.exe",
	// 	Arguments: []string{
	// 		"P:/projects/service-wrapper/test.ps1",
	// 	},
	// 	LogDirectory: "P:/projects/service-wrapper/log",
	// }, isAnInteractiveSession)
}
