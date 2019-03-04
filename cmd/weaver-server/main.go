package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	raven "github.com/getsentry/raven-go"
	"github.com/gojektech/weaver/config"
	"github.com/gojektech/weaver/etcd"
	"github.com/gojektech/weaver/pkg/instrumentation"
	"github.com/gojektech/weaver/pkg/logger"
	"github.com/gojektech/weaver/server"
	cli "gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "weaver"
	app.Usage = "run weaver-server"
	app.Version = fmt.Sprintf("%s built on %s (commit: %s)", Version, BuildDate, Commit)
	app.Description = "An Advanced HTTP Reverse Proxy with Dynamic Sharding Strategies"
	app.Commands = []cli.Command{
		{
			Name:        "start",
			Description: "Start weaver server",
			Action:      startWeaver,
		},
		{
			Name:        "acls",
			Aliases:     []string{"a"},
			Description: "List, Create, Delete, Update ACLs",
			Usage:       "Perform list, create, update, delete acls",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "etcd-host, etcd",
					Value:  "http://localhost:2379",
					Usage:  "Host address of ETCD",
					EnvVar: "ETCD_ENDPOINTS",
				},
				cli.StringFlag{
					Name:   "namespace, ns",
					Value:  "weaver",
					Usage:  "Namespace of Weaver ACLS",
					EnvVar: "ETCD_KEY_PREFIX",
				},
			},
			Subcommands: []cli.Command{
				{
					Name:    "list",
					Usage:   "List ACLS",
					Aliases: []string{"l"},
					Action: func(c *cli.Context) error {
						setupEnv(c)
						config.Load()
						rl, err := etcd.NewRouteLoader()
						if err != nil {
							return err
						}
						acls, err := rl.ListAll()
						if err != nil {
							return err
						}
						fmt.Println(acls)
						return nil
					},
				},
				{
					Name:    "create",
					Usage:   "Create ACL",
					Aliases: []string{"c"},
					Action: func(c *cli.Context) error {
						setupEnv(c)
						fmt.Println("new task template: ", c.Args().First())
						return nil
					},
				},
				{
					Name:    "update",
					Usage:   "Update ACL",
					Aliases: []string{"u"},
					Action: func(c *cli.Context) error {
						setupEnv(c)
						fmt.Println("new task template: ", c.Args().First())
						return nil
					},
				},
				{
					Name:    "delete",
					Usage:   "Delete ACL",
					Aliases: []string{"d"},
					Action: func(c *cli.Context) error {
						setupEnv(c)
						fmt.Println("new task template: ", c.Args().First())
						return nil
					},
				},
			},
		},
	}

	app.Run(os.Args)
}

func setupEnv(c *cli.Context) {
	os.Setenv("ETCD_ENDPOINTS", c.GlobalString("etcd-host"))
	os.Setenv("ETCD_KEY_PREFIX", c.GlobalString("namespace"))
	config.Load()
}

func startWeaver(_ *cli.Context) error {
	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	config.Load()

	raven.SetDSN(config.SentryDSN())
	logger.SetupLogger()

	err := instrumentation.InitiateStatsDMetrics()
	if err != nil {
		log.Printf("StatsD: Error initiating client %s", err)
	}
	defer instrumentation.CloseStatsDClient()

	instrumentation.InitNewRelic()
	defer instrumentation.ShutdownNewRelic()

	routeLoader, err := etcd.NewRouteLoader()
	if err != nil {
		log.Printf("StartServer: failed to initialise etcd route loader: %s", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go server.StartServer(ctx, routeLoader)

	sig := <-sigC
	log.Printf("Received %d, shutting down", sig)

	defer cancel()
	server.ShutdownServer(ctx)

	return nil
}

// Build information (will be injected during build)
var (
	Version   = "1.0.0"
	Commit    = "n/a"
	BuildDate = "n/a"
)
