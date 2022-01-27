package aws

import (
	"github.com/RaniSputnik/events"
)

type ARN = events.ID

// TODO: Use these aliases when constructing ARN's

var aliases map[string]string

func Aliases(names map[string]string) {
	if names == nil {
		names = map[string]string{}
	}
	aliases = names
}
