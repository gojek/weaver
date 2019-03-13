package cli

import (
	"fmt"
	"github.com/gojektech/weaver/config"
	"github.com/gojektech/weaver/pkg/logger"
	cliParser "gopkg.in/urfave/cli.v1"
	"os"
)

const (
	aclCliName     = "acls"
	aclDescription = "List, Create, Delete, Update ACLs"
	aclUsage       = "Perform list, create, update, delete acls"
)

var (
	aclCommandMap = make(map[string]Command)
	aclAliases    = []string{"a"}
	aclFlags      = []cliParser.Flag{
		cliParser.StringFlag{
			Name:   "etcd-host, etcd",
			Value:  "http://localhost:2379",
			Usage:  "Host address of ETCD",
			EnvVar: "ETCD_ENDPOINTS",
		},
		cliParser.StringFlag{
			Name:   "namespace, ns",
			Value:  "weaver",
			Usage:  "Namespace of Weaver ACLS",
			EnvVar: "ETCD_KEY_PREFIX",
		},
	}

	aclCommand = &struct {
		flags []cliParser.Flag
		*command
	}{
		command: &command{
			cliName:     aclCliName,
			description: aclDescription,
			usage:       aclUsage,
			aliases:     aclAliases,
			getCommand: func() cliParser.Command {
				subCommands := populateSubCommands(&aclCommandMap)
				cliParserCmd := cliParser.Command{
					Name:        string(aclCliName),
					Description: aclDescription,
					Usage:       aclUsage,
					Aliases:     aclAliases,
					Flags:       aclFlags,
				}
				if len(subCommands) > 0 {
					cliParserCmd.Subcommands = subCommands
				} else {
					cliParserCmd.Action = ExecCommand
				}
				return cliParserCmd
			},
			action: func(c *cliParser.Context) error {
				preSetupForACLCli(c)
				if cmdToExecute := aclCommandMap[c.Command.Name]; cmdToExecute != nil {
					return cmdToExecute.Exec(c)
				}
				return fmt.Errorf("Error executing command. Command not registered: %s", c.Command.Name)
			},
		},
		flags: aclFlags,
	}
)

func registerACLCommand(cliName string, cmd Command) {
	if cmdToExecute := aclCommandMap[cliName]; cmdToExecute != nil {
		panic(fmt.Sprintf("Command for cli %s already registered by: %s", cliName, cmdToExecute))
	}
	aclCommandMap[cliName] = cmd
}

func populateSubCommands(cmdMap *map[string]Command) []cliParser.Command {
	subCommands := []cliParser.Command{}
	for _, registerCmd := range *cmdMap {
		subCommand := registerCmd.GetCommand()
		subCommands = append(subCommands, subCommand)
	}
	return subCommands
}

func preSetupForACLCli(c *cliParser.Context) {
	os.Setenv("ETCD_ENDPOINTS", c.GlobalString("etcd-host"))
	os.Setenv("ETCD_KEY_PREFIX", c.GlobalString("namespace"))
	config.Load()
	logger.SetupLogger()
}

func init() {
	registerCommand(aclCliName, aclCommand)
}
