package events

import (
	"context"
)

type Event struct {
	Name        string   `json:"name"`
	Publishers  []string `json:"publishers"`
	Subscribers []string `json:"subscribers"`
}

type Service interface {
	Get(ctx context.Context, eventName string) (*Event, error)
}
