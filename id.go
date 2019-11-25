package events

import "fmt"

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

func (gid GID) String() string {
	return fmt.Sprintf("gid:%s:%s:%s", gid.AccountID, gid.Kind, gid.ID)
}
