package aws

import (
	"context"
	"encoding/json"

	"github.com/RaniSputnik/events"
	"github.com/aws/aws-sdk-go/aws"
	cw "github.com/aws/aws-sdk-go/service/cloudwatchevents"
)

// getAllSubscribers fetches the GIDs of all of the infrastructure
// subscribing to an event in AWS with the given name. Returns either
// the subscribers or an error.
func (s *service) getAllSubscribers(ctx context.Context, eventName string) ([]events.GID, error) {
	// When this function exits, this context will be
	// cancelled, stopping all routines in the pipeline
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// The event that we are expecting to be published
	// from Lambda functions - with correct source and
	// detail type (AKA namespace & event name)
	event := Event{
		Source:     s.namespace,
		DetailType: eventName,
	}

	// Create an error channel for API errors from AWS
	// It's possible that both functions return an error
	// so we buffer to ensure that both can write without
	// blocking in this scenario.
	errc := make(chan error, 2)

	// Create the pipeline for rules processing
	// We start by streaming all CW rules in the account
	// which we then filter to only rules that would match
	// the given event. We then fetch all targets for each
	// of these rules and convert those targets into
	// "Subscribers" (our domain model)
	allRulesc := s.rulesStream(ctx, errc)
	matchingRulesc := s.filterRulesStream(ctx, allRulesc, ruleWouldMatchEventFilter(event))
	targetsc := s.targetsStream(ctx, matchingRulesc, errc)
	subc := s.subscriberStream(ctx, targetsc)

	subscribers := []events.GID{}
	for {
		select {
		// Drain all of the subscribers from the pipeline.
		case gid := <-subc:
			if gid.ID == "" {
				// There are no
				return subscribers, nil
			}
			subscribers = append(subscribers, gid)
		// Alternatively, exit as soon as an error is encountered.
		case err := <-errc:
			return nil, err
		// Alternatively, alternatively, cancel work if the context
		// has been cancelled or has timed out.
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func (s *service) rulesStream(ctx context.Context, errc chan<- error) <-chan *cw.Rule {
	pageSize := 10
	out := make(chan *cw.Rule, pageSize)

	go func() {
		// Deferring ensures that whenever we exit this
		// function, the output channel will be closed
		defer close(out)

		var nextToken *string
		for {
			// Fetch CloudWatch rules using the AWS SDK
			params := cw.ListRulesInput{
				Limit:     aws.Int64(int64(pageSize)),
				NextToken: nextToken,
			}
			Debugf("cloudwatchevents.ListRules: %+v\n", params)
			output, err := s.client.ListRules(&params)
			// Return an error on the error channel if
			// there was one, then abandon further execution
			if err != nil {
				Debugf("cloudwatchevents.ListRules failed: %+v", err)
				errc <- ctx.Err()
				return
			}

			// We either send a new result down the
			// output channel OR if the context has been
			// cancelled then we abandon execution
			for _, result := range output.Rules {
				select {
				case out <- result:
				case <-ctx.Done():
					return
				}
			}

			// If there is no next page, we are done
			nextToken = output.NextToken
			if nextToken == nil {
				return
			}
		}
	}()

	return out
}

// filterRulesStream takes a stream of rules as input and filters
// the stream using the ruleFilter given. It returns a new channel
// of rules that will contain only the rules that match the given
// rule filter.
// The rule filter should return false to exclude a given rule from
// being sent to the output stream.
func (s *service) filterRulesStream(ctx context.Context, in <-chan *cw.Rule, allow ruleFilter) <-chan *cw.Rule {
	out := make(chan *cw.Rule)
	go func() {
		defer close(out)
		for rule := range in {
			// Ignore rules that don't match the filter
			if !allow(rule) {
				continue
			}
			select {
			// Push a matching rule to the output stream
			case out <- rule:
			// Abandon work if the context has been cancelled
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

func (s *service) targetsStream(ctx context.Context, in <-chan *cw.Rule, errc chan<- error) <-chan *cw.Target {
	out := make(chan *cw.Target, 10) // TODO: How big of a buffer to use?
	go func() {
		defer close(out)
		for rule := range in {
			input := cw.ListTargetsByRuleInput{
				Rule: rule.Name,
				// TODO: Instead, handle pagination here
				Limit: aws.Int64(maxTargetsPageSize),
			}
			Debugf("cloudwatchevents.ListTargetsByRule: %+v\n", input)
			desc, err := s.client.ListTargetsByRule(&input)
			if err != nil {
				Debugf("cloudwatchevents.ListTargetsByRule failed: %+v", err)
				errc <- err
				return
			}

			for _, target := range desc.Targets {
				out <- target
			}
		}
	}()
	return out
}

func (s *service) subscriberStream(ctx context.Context, in <-chan *cw.Target) <-chan events.GID {
	out := make(chan events.GID)
	go func() {
		defer close(out)
		for target := range in {
			arn := aws.StringValue(target.Arn)
			gid := ARNToGID(arn)
			// Special case to ensure we set the namespace GID correctly
			if gid.Kind == events.KindNamespace && gid.ID == "default" {
				gid.ID = s.namespace
			}
			out <- gid
		}
	}()
	return out
}

type ruleFilter func(*cw.Rule) bool

func ruleWouldMatchEventFilter(event Event) ruleFilter {
	return func(r *cw.Rule) bool {
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
