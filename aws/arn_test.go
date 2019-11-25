package aws_test

import (
	"testing"

	"github.com/RaniSputnik/events"
	"github.com/RaniSputnik/events/aws"
	"github.com/stretchr/testify/assert"
)

func TestARNToGID(t *testing.T) {
	tests := []struct {
		Given  string
		Expect events.GID
	}{
		{
			"arn:aws:lambda:eu-west-1:123456789012:function:someFunction",
			events.GID{AccountID: "123456789012", Kind: events.KindFunction, ID: "someFunction"},
		},
		{
			"arn:aws:lambda:eu-west-1:999999999999:function:someFunction",
			events.GID{AccountID: "999999999999", Kind: events.KindFunction, ID: "someFunction"},
		},
		{
			"arn:aws:events:eu-west-1:513286705436:event-bus/default",
			events.GID{AccountID: "513286705436", Kind: events.KindNamespace, ID: "default"},
		},
	}

	for _, test := range tests {
		t.Run(test.Given, func(t *testing.T) {
			assert.Equal(t, test.Expect, aws.ARNToGID(test.Given))
		})
	}
}
