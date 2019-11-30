package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchevents"
)

func (s *service) rulesStream(ctx context.Context, errc chan<- error) <-chan *cloudwatchevents.Rule {
	pageSize := 10
	results := make(chan *cloudwatchevents.Rule, pageSize)

	go func(results chan<- *cloudwatchevents.Rule) {
		defer close(results)

		var nextToken *string
		for {
			params := cloudwatchevents.ListRulesInput{
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
					errc <- ctx.Err()
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
