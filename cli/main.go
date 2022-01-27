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
	source := flag.String("source", "", "The source to query for events")
	eventName := flag.String("event", "", "The name of the event you want to query")
	verbose := flag.Bool("v", false, "Verbose mode, optionally prints debug logs")
	flag.Parse()

	if *verbose {
		aws.Debugf = log.Printf // Print debug logs to the console
	}

	if *source == "" {
		fmt.Println("Missing required argument 'source'")
		flag.Usage()
		return
	}
	if *eventName == "" {
		fmt.Println("Missing required argument 'event'")
		flag.Usage()
		return
	}

	events := aws.Events(*source)
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
