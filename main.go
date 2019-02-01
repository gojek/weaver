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
	"github.com/gojekfarm/weaver/internal/config"
	"github.com/gojekfarm/weaver/internal/server"
	"github.com/gojekfarm/weaver/pkg/instrumentation"
	"github.com/gojekfarm/weaver/pkg/logger"
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

	go server.StartServer()

	sig := <-sigC
	log.Printf("Received %d, shutting down", sig)

	ctx, cancel := context.WithTimeout(context.Background(), (1 * time.Second))
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
