package cli_test

import (
	"fmt"
	"github.com/gojektech/weaver/internal/cli"
	"github.com/stretchr/testify/assert"
	baseCli "gopkg.in/urfave/cli.v1"
	"testing"
)

type testFlagSetup struct {
	name, value, usage, env string
	flag                    baseCli.Flag
}

func TestInitializationOfStringFlag(t *testing.T) {
	ts := setupTestFlag()

	assert.Equal(t, ts.flag.String(), fmt.Sprintf("--%s value\t%s (default: \"%s\") [$%s]", ts.name, ts.usage, ts.value, ts.env))
}

func setupTestFlag() *testFlagSetup {
	setup := &testFlagSetup{name: "test", value: "default", usage: "use to set value for test", env: "TEST"}
	setup.flag = cli.NewStringFlag(setup.name, setup.value, setup.usage, setup.env)
	return setup
}
