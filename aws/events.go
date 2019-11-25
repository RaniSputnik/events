package aws

import (
	"context"

	"encoding/json"

	"github.com/RaniSputnik/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchevents"
)

const maxRulesPageSize = 100
const maxTargetsPageSize = 100

type service struct {
	namespace string
	client    *cloudwatchevents.CloudWatchEvents
}

func Events(namespace string) events.Service {
	s := session.Must(session.NewSession())
	client := cloudwatchevents.New(s)
	return &service{
		namespace: namespace,
		client:    client,
	}
}

func (s *service) Get(ctx context.Context, eventName string) (res *events.Event, err error) {
	params := cloudwatchevents.ListRulesInput{
		Limit: aws.Int64(maxRulesPageSize),
	}
	var output *cloudwatchevents.ListRulesOutput
	output, err = s.client.ListRules(&params)
	if err != nil || output == nil || len(output.Rules) == 0 {
		return
	}
	rulesForEvent := applyRuleFilter(eventNamedFilter(s.namespace, eventName), output.Rules)

	subscribers := []string{}
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
			subscribers = append(subscribers, gid.String())
		}
	}

	return &events.Event{
		Name:        eventName,
		Subscribers: subscribers,
		Publishers:  []string{},
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

type eventPattern struct {
	Source     []string `json:"source"`
	Account    []string `json:"account"`
	DetailType []string `json:"detail-type"`
}

func eventNamedFilter(namespace string, name string) ruleFilter {
	return func(r *cloudwatchevents.Rule) bool {
		if r.EventPattern == nil {
			return false
		}
		pattern := []byte(aws.StringValue(r.EventPattern))
		var decoded eventPattern
		if err := json.Unmarshal(pattern, &decoded); err != nil {
			return false
		}
		if len(decoded.Source) > 0 {
			found := false
			for _, source := range decoded.Source {
				if source == namespace {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		for _, receiveEvent := range decoded.DetailType {
			if receiveEvent == name {
				return true
			}
		}
		return false
	}
}
