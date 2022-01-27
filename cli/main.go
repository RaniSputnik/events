package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/RaniSputnik/events/aws"
)

func main() {
	namespace := flag.String("namespace", "", "The namespace to query for events")
	eventName := flag.String("event", "", "The name of the event you want to query")
	verbose := flag.Bool("v", false, "Verbose mode, optionally prints debug logs")
	flag.Parse()

	if *verbose {
		aws.Debugf = log.Printf // Print debug logs to the console
	}

	if *namespace == "" {
		fmt.Println("Missing required argument 'namespace'")
		flag.Usage()
		return
	}
	if *eventName == "" {
		fmt.Println("Missing required argument 'event'")
		flag.Usage()
		return
	}

	events := aws.Events(*namespace)
	ev, err := events.Get(context.Background(), *eventName)
	if err != nil {
		exit(err)
	}
	for _, sub := range ev.Subscribers {
		fmt.Println(sub)
	}
}

func exit(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
