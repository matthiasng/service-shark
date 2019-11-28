package main

import (
	"flag"
	"fmt"
	"github.com/bioapfelsaft/service-wrapper/service"
	"golang.org/x/sys/windows/svc"
	"log"
)

func main() {
	serviceName := flag.String("name", "", "ServiceName")

	flag.Parse()

	fmt.Println(*serviceName)

	isAnInteractiveSession, err := svc.IsAnInteractiveSession()	 // #todo Wird das von der UI Interaction Einstellung beeinflusst ?
	if err != nil {
		log.Fatalf("failed to determine if we are running in an interactive session: %v", err)
	}

	service.Run(*serviceName, isAnInteractiveSession)
	return
}
