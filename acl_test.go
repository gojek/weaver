package weaver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	validACLDoc = `{
		"id": "gojek_hello",
		"criterion" : "Method('POST') && Path('/gojek/hello-service')",
		"endpoint" : {
			"shard_expr": ".serviceType",
			"matcher": "body",
			"shard_func": "lookup",
			"shard_config": {
				"999": {
					"backend_name": "hello_backend",
					"backend":"http://hello.golabs.io"
				}
			}
		}
	}`
)

var genACLTests = []struct {
	aclDoc      string
	expectedErr error
}{
	{
		validACLDoc,
		nil,
	},
}

func TestGenACL(t *testing.T) {
	acl := &ACL{}

	for _, tt := range genACLTests {
		actualErr := acl.GenACL(tt.aclDoc)
		if actualErr != tt.expectedErr {
			assert.Failf(t, "Unexpected error message", "want: %v got: %v", tt.expectedErr, actualErr)
		}
	}
}
