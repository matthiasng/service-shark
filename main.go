package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bioapfelsaft/service-wrapper/service"
	"golang.org/x/sys/windows/svc"
)

// #todo config ?
// 1. -bind:NAME -> -NAME <Value>
// 2. -NAME cfg:VALUE_NAME -> -NAME <Value>

// wrapper.exe -cfg-file test.yaml
// $env:NAME="123" + binding

// load order
// 1. config file
// 2. environment
// 3. wrapper.exe -arguments "..."

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

	//
	args, err := parseCommandLine(*arguments)
	if err != nil {
		log.Fatalf("invalid arguments: %v", err)
	}

	//
	for i := 0; i < len(args); i++ {
		if strings.HasPrefix(args[i], "bind:") {
			args[i] = os.Getenv(strings.TrimLeft(args[i], "bind:"))
		}
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
}
