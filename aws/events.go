package aws

import (
	"context"

	"github.com/RaniSputnik/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eventbridge"
	"github.com/aws/aws-sdk-go/service/sts"
)

var Debugf = func(format string, v ...interface{}) { /* No op */ }

const maxRulesPageSize = 100
const maxTargetsPageSize = 100

type service struct {
	namespace string
	client    *eventbridge.EventBridge
	sts       *sts.STS
}

func Events(namespace string) events.Service {
	s := session.Must(session.NewSession())
	client := eventbridge.New(s)
	sts := sts.New(s)
	return &service{
		namespace: namespace,
		client:    client,
		sts:       sts,
	}
}

func (s *service) Get(ctx context.Context, eventName string) (res *events.Event, err error) {
	var caller *sts.GetCallerIdentityOutput
	caller, err = s.sts.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	currentAccountID := aws.StringValue(caller.Account)
	if alias, hasAlias := aliases[currentAccountID]; hasAlias {
		currentAccountID = alias
	}

	var subscribers []events.GID
	if subscribers, err = s.getAllSubscribers(ctx, eventName); err != nil {
		return
	}

	// Create the resulting event
	return &events.Event{
		ID: eventName,
		GID: events.GID{
			AccountID: currentAccountID,
			Kind:      events.KindSchema,
			ID:        eventName,
		},
		Name:        eventName,
		Subscribers: subscribers,
		Publishers: []events.EventPublisher{
			{
				Function:  events.GID{AccountID: currentAccountID, ID: "todo", Kind: events.KindFunction},
				Namespace: events.GID{AccountID: currentAccountID, ID: "todo", Kind: events.KindNamespace},
			},
		},
	}, nil
}
