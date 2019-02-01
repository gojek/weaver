package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	app.Version = fmt.Sprintf("%s built on %s (commit: %s)", Version, BuildDate, Commit)
	app.Description = "A Layer-7 Load Balancer with Dynamic Sharding Strategies"
	app.Commands = []cli.Command{
		{
			Name:        "server",
			Description: "Start weaver server",
			Action:      startWeaver,
		},
	}

	app.Run(os.Args)
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

	ctx, cancel := context.WithTimeout(context.Background(), (1 * time.Second))
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
