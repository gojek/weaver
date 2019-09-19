package commands

import (
	"github.com/gojektech/weaver/internal/cli"
)

var ETCDFLAG = cli.NewStringFlag("etcd-host, etcd", "http://localhost:2379", "HOST Address of ETCD", "ETCD_ENDPOINTS")
var NAMESPACEFLAG = cli.NewStringFlag("namespace, ns", "weaver", "Namespace of Weaver ACLS", "ETCD_KEY_PREFIX")
