package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/RaniSputnik/events/aws"
)

func main() {
	namespace := flag.String("namespace", "", "The namespace to query for events")
	eventName := flag.String("event", "", "The name of the event you want to query")
	flag.Parse()

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
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "\t")
	if err = encoder.Encode(ev); err != nil {
		exit(err)
	}
}

func exit(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
