package events

import (
	"errors"
	"fmt"
	"strings"
)

type kind string

const (
	KindFunction  = kind("function")
	KindSchema    = kind("schema")
	KindNamespace = kind("namespace")
)

type GID struct {
	AccountID string
	Kind      kind
	ID        string
}

func ParseGID(str string) (GID, error) {
	parts := strings.Split(str, ":")
	if len(parts) != 4 || parts[0] != "gid" {
		return GID{}, errors.New("invalid gid: '" + str + "'")
	}
	return GID{
		AccountID: parts[1],
		Kind:      kind(parts[2]), // TODO: Validate me
		ID:        parts[3],
	}, nil
}

func (gid GID) String() string {
	return fmt.Sprintf("gid:%s:%s:%s", gid.AccountID, gid.Kind, gid.ID)
}

func (a *GID) UnmarshalJSON(b []byte) error {
	parsed, err := ParseGID(string(b))
	if err != nil {
		return err
	}
	a.AccountID = parsed.AccountID
	a.Kind = parsed.Kind
	a.ID = parsed.ID
	return nil
}

func (a GID) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, a)), nil
}
