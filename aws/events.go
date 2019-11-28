package aws

import (
	"context"

	"encoding/json"

	"github.com/RaniSputnik/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchevents"
	"github.com/aws/aws-sdk-go/service/sts"
)

const maxRulesPageSize = 100
const maxTargetsPageSize = 100

type service struct {
	namespace string
	client    *cloudwatchevents.CloudWatchEvents
	sts       *sts.STS
}

func Events(namespace string) events.Service {
	s := session.Must(session.NewSession())
	client := cloudwatchevents.New(s)
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

	params := cloudwatchevents.ListRulesInput{
		Limit: aws.Int64(maxRulesPageSize),
	}
	var output *cloudwatchevents.ListRulesOutput
	output, err = s.client.ListRules(&params)
	if err != nil || output == nil || len(output.Rules) == 0 {
		return
	}
	event := Event{
		Source:     s.namespace,
		DetailType: eventName,
	}
	rulesForEvent := applyRuleFilter(eventMatchFilter(event), output.Rules)

	subscribers := []events.GID{}
	for _, rule := range rulesForEvent {
		// TODO: run this in parallel
		input := cloudwatchevents.ListTargetsByRuleInput{
			Rule:  rule.Name,
			Limit: aws.Int64(maxTargetsPageSize),
		}
		desc, err := s.client.ListTargetsByRule(&input)
		if err != nil {
			return nil, err
		}
		for _, target := range desc.Targets {
			arn := aws.StringValue(target.Arn)
			gid := ARNToGID(arn)
			// Special case to ensure we set the namespace GID correctly
			if gid.Kind == events.KindNamespace && gid.ID == "default" {
				gid.ID = s.namespace
			}
			subscribers = append(subscribers, gid)
		}
	}

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

type ruleFilter func(*cloudwatchevents.Rule) bool

func applyRuleFilter(f ruleFilter, rules []*cloudwatchevents.Rule) []*cloudwatchevents.Rule {
	results := make([]*cloudwatchevents.Rule, 0, len(rules))
	for _, rule := range rules {
		if f(rule) {
			results = append(results, rule)
		}
	}
	return results
}

func eventMatchFilter(event Event) ruleFilter {
	return func(r *cloudwatchevents.Rule) bool {
		if r.EventPattern == nil {
			return false // TODO: Does a nil pattern mean no match? Or match everything?
		}
		pattern := []byte(aws.StringValue(r.EventPattern))
		var decoded EventPattern
		if err := json.Unmarshal(pattern, &decoded); err != nil {
			return false
		}
		return EventMatches(event, decoded)
	}
}
