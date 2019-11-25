package aws_test

import (
	"testing"

	"github.com/RaniSputnik/events/aws"
	"github.com/stretchr/testify/assert"
)

func TestEventMatches(t *testing.T) {
	anyEvent := aws.Event{
		Source:     "com.example",
		DetailType: "someEvent",
	}

	t.Run("EmptyPatternMatchesAnything", func(t *testing.T) {
		assert.True(t, aws.EventMatches(anyEvent, aws.EventPattern{}))
	})

	t.Run("DoesNotMatchWithIncorrectSource", func(t *testing.T) {
		event := aws.Event{
			Source:     "com.example",
			DetailType: "someEvent",
		}
		pattern := aws.EventPattern{
			Source: []string{"uk.co.example"},
		}
		assert.False(t, aws.EventMatches(event, pattern))
	})

	t.Run("MatchesIfSourcePresentInList", func(t *testing.T) {
		event := aws.Event{
			Source:     "com.example",
			DetailType: "someEvent",
		}
		pattern := aws.EventPattern{
			Source: []string{"uk.co.example", "com.example"},
		}
		assert.True(t, aws.EventMatches(event, pattern))
	})

	t.Run("DoesNotMatchWithIncorrectDetailType", func(t *testing.T) {
		event := aws.Event{
			Source:     "com.example",
			DetailType: "someEvent",
		}
		pattern := aws.EventPattern{
			DetailType: []string{"anotherEvent"},
		}
		assert.False(t, aws.EventMatches(event, pattern))
	})

	t.Run("MatchesIfDetailTypePresentInList", func(t *testing.T) {
		event := aws.Event{
			Source:     "com.example",
			DetailType: "someEvent",
		}
		pattern := aws.EventPattern{
			DetailType: []string{"anotherEvent", "someEvent"},
		}
		assert.True(t, aws.EventMatches(event, pattern))
	})

	t.Run("DoesNotMatchIfSourceValidButDetailTypeInvalid", func(t *testing.T) {
		event := aws.Event{
			Source:     "com.example",
			DetailType: "someEvent",
		}
		pattern := aws.EventPattern{
			Source:     []string{"com.example"},
			DetailType: []string{"anotherEvent"},
		}
		assert.False(t, aws.EventMatches(event, pattern))
	})

	t.Run("DoesNotMatchIfDetailTypeValidButSourceInvalid", func(t *testing.T) {
		event := aws.Event{
			Source:     "com.example",
			DetailType: "someEvent",
		}
		pattern := aws.EventPattern{
			Source:     []string{"uk.co.example"},
			DetailType: []string{"someEvent"},
		}
		assert.False(t, aws.EventMatches(event, pattern))
	})
}
