package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	cw "github.com/aws/aws-sdk-go/service/cloudwatchevents"
)

func (s *service) rulesStream(ctx context.Context, errc chan<- error) <-chan *cw.Rule {
	pageSize := 10
	results := make(chan *cw.Rule, pageSize)

	go func(results chan<- *cw.Rule) {
		defer close(results)

		var nextToken *string
		for {
			params := cw.ListRulesInput{
				Limit:     aws.Int64(int64(pageSize)),
				NextToken: nextToken,
			}
			output, err := s.client.ListRules(&params)
			if err != nil {
				errc <- ctx.Err()
				return
			}
			for _, result := range output.Rules {
				select {
				case results <- result:
				case <-ctx.Done():
					return
				}
			}
			nextToken = output.NextToken
			if nextToken == nil {
				return
			}
		}
	}(results)

	return results
}

func (s *service) filterRulesStream(ctx context.Context, in <-chan *cw.Rule, allow ruleFilter) <-chan *cw.Rule {
	out := make(chan *cw.Rule)
	go func() {
		defer close(out)
		for rule := range in {
			if !allow(rule) {
				continue
			}
			select {
			case out <- rule:
			case <-ctx.Done():
			}
		}
	}()
	return out
}
