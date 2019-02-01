package server

import (
	"encoding/json"
	"fmt"
)

// ACL - Connects to an external endpoint
type ACL struct {
	ID             string          `json:"id"`
	Criterion      string          `json:"criterion"`
	EndpointConfig *EndpointConfig `json:"endpoint"`

	Endpoint *Endpoint
}

// GenACL - Generates an ACL from JSON
func (acl *ACL) GenACL(val string) error {
	return json.Unmarshal([]byte(val), &acl)
}

func (acl ACL) String() string {
	return fmt.Sprintf("ACL(%s, %s)", acl.ID, acl.Criterion)
}
