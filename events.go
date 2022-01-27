package events

import (
	"context"
)

type ID string

type EventPublisher struct {
	Function  ID `json:"function"`
	Namespace ID `json:"namespace"`
}

type Event struct {
	ID          ID           	 `json:"id"`
	Name        string           `json:"name"`
	Publishers  []EventPublisher `json:"publishers"`
	Subscribers []ID             `json:"subscribers"`
}

type Service interface {
	Get(ctx context.Context, eventName string) (*Event, error)
}
