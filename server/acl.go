package server

import (
	"encoding/json"
	"fmt"

	"github.com/gojektech/weaver"
)

// ACL - Connects to an external endpoint
type ACL struct {
	ID             string                 `json:"id"`
	Criterion      string                 `json:"criterion"`
	EndpointConfig *weaver.EndpointConfig `json:"endpoint"`

	Endpoint *weaver.Endpoint
}

// GenACL - Generates an ACL from JSON
func (acl *ACL) GenACL(val string) error {
	return json.Unmarshal([]byte(val), &acl)
}

func (acl ACL) String() string {
	return fmt.Sprintf("ACL(%s, %s)", acl.ID, acl.Criterion)
}
