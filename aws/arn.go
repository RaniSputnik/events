package aws

import (
	"strings"

	"github.com/RaniSputnik/events"
)

func ARNToGID(arn string) (res events.GID) {
	failureResponse := events.GID{}
	parts := strings.Split(arn, ":")
	res.AccountID = parts[4]
	resourceType := parts[5]
	switch resourceType {
	case "function":
		res.Kind = events.KindFunction
		res.ID = parts[6]
	case "event-bus/default":
		res.Kind = events.KindNamespace
		res.ID = "default"
	default:
		return failureResponse // Unknown
	}
	return
}
