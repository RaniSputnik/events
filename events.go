package events

import (
	"context"
)

type EventPublisher struct {
	Function  GID `json:"function"`
	Namespace GID `json:"namespace"`
}

type Event struct {
	ID          string           `json:"id"`
	GID         GID              `json:"gid"` // TODO: Derive this field from ID?
	Name        string           `json:"name"`
	Publishers  []EventPublisher `json:"publishers"`
	Subscribers []GID            `json:"subscribers"`
}

type Service interface {
	Get(ctx context.Context, eventName string) (*Event, error)
}
