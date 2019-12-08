package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/matthiasng/service-shark/cli"
	"github.com/matthiasng/service-shark/command"
	"github.com/matthiasng/service-shark/service"
)

var arguments = cli.Arguments{}

func init() {
	flag.CommandLine.Usage = func() {
		fmt.Println("Usage of Service Shark:")
		flag.PrintDefaults()
	}

	flag.StringVar(&arguments.Name, "name", "", "Service name [required]")
	flag.StringVar(&arguments.WorkingDirectory, "workdir", "", "Working directory [required]")
	flag.StringVar(&arguments.LogDirectory, "logdir", "./log", "Log directory. Absolute path or relative to working directory.")
	flag.StringVar(&arguments.Command, "cmd", "", "Command [required]")

	flag.Parse()
}

func main() {
	if err := service.FixWorkingDirectory(); err != nil {
		log.Fatal(err)
	}

	arguments.CommandArguments = cli.ExpandArguments(flag.Args())

	if err := cli.Validate(arguments); err != nil {
		log.Fatal(err)
	}

	if err := os.Chdir(arguments.WorkingDirectory); err != nil {
		log.Fatal(err)
	}

	host := command.Host{
		Arguments: arguments,
	}

	if err := service.Run(&host); err != nil {
		log.Fatal(err)
	}
}
