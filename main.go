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

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var (
	printVersion *bool = nil
	arguments          = cli.Arguments{}
)

func init() {
	flag.CommandLine.Usage = func() {
		fmt.Println("Service Shark:")
		fmt.Println("  Host any executable as a Windows service.")
		fmt.Println("Usage:")
		flag.PrintDefaults()
		fmt.Println("  -- (terminator)")
		fmt.Println(`	Pass all arguments after the terminator "--" to command.`)
		fmt.Println(`	Bind argument to environment variable with "env:{VAR_NAME}".`)
		fmt.Println("Example:")
		fmt.Println(`  service-shark.exe ... -cmd java -- -jar test.jar -Xmx1G -myArg "env:MY_VALUE"`)
		fmt.Println(`  => java -jar test.jar -Xmx1G -myArg "my env value"`)
	}

	flag.StringVar(&arguments.Name, "name", "", "Service name [required]")
	flag.StringVar(&arguments.WorkingDirectory, "workdir", "./", "Working directory")
	flag.StringVar(&arguments.LogDirectory, "logdir", "./log", "Log directory.\nFile name: {name}_YYYY-MM-DD_HH-MM-SS")
	flag.StringVar(&arguments.Command, "cmd", "", `Command [required]`)
	printVersion = flag.Bool("version", false, "Print version and exit")

	flag.Parse()
}

func main() {
	if *printVersion {
		fmt.Printf("%v, commit %v, built at %v", version, commit, date)
		os.Exit(0)
	}

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
