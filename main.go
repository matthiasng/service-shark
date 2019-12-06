package main

import (
	"log"

	"github.com/matthiasng/service-shark/command"
	"github.com/matthiasng/service-shark/service"
)

// #todo working directory

// 	parser := argparse.NewParser("service-wrapper", `Run a "-command" with "-arguments" as service`)

// 	serviceName := parser.String("n", "name", &argparse.Options{
// 		Required: true,
// 		Help:     "Servicename",
// 	})
// 	logDirectory := parser.String("l", "logdirectory", &argparse.Options{
// 		Required: true,
// 		Help:     "Log directory",
// 	})
// 	command := parser.String("c", "command", &argparse.Options{
// 		Required: true,
// 		Help:     "Command",
// 	})
// 	arguments := parser.List("a", "arg", &argparse.Options{
// 		Required: true,
// 		Help:     `Command arguments. Example: '... -a "-key" -a "value" -a "--key2" -a "value"'`,
// 	})

// 	err := parser.Parse(os.Args)
// 	if err != nil {
// 		fmt.Print(parser.Usage(err))
// 		os.Exit(2)
// 	}

// cli.BindArguments(*arguments),

func main() {
	host := command.Host{
		CmdConfig: command.Config{
			//Name: "powershell",
			Name: "C:/Program Files/PowerShell/7-preview/preview/pwsh-preview.cmd",
			Arguments: []string{
				`P:\_dev\projects\service-shark\example\test-service.ps1`,
			},
		},
		LogDirecotry: "c:/tmp",
	}

	if err := service.Run(&host, "todo"); err != nil {
		log.Fatalf("[error] - %v", err)
	}
}
