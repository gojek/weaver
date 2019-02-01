package etcd

import (
	"fmt"

	"github.com/gojektech/weaver"
)

const (
	// ACLKeyFormat - Format for a ACL's key in a KV Store
	ACLKeyFormat = "/%s/acls/%s/acl"
)

// ACLKey - Points to a stored ACL
type ACLKey string

// GenACLKey - Generate an ACL Key given etcd's node key
func GenACLKey(key string) ACLKey {
	return ACLKey(fmt.Sprintf("%s/acl", key))
}

func GenKey(acl *weaver.ACL, pfx string) ACLKey {
	return ACLKey(fmt.Sprintf(ACLKeyFormat, pfx, acl.ID))
}
