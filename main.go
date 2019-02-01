package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gojektech/weaver/internal/config"
	"github.com/gojektech/weaver/internal/server"
	"github.com/gojektech/weaver/pkg/instrumentation"
	"github.com/gojektech/weaver/pkg/logger"

	raven "github.com/getsentry/raven-go"
	cli "gopkg.in/urfave/cli.v1"
)

func version(_ *cli.Context) error {
	fmt.Println(GenVersion())
	return nil
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

func main() {
	app := cli.NewApp()
	app.Name = "server"
	app.Description = "A Layer-7 Load Balancer with Dynamic Sharding Strategies"
	app.Commands = []cli.Command{
		{
			Name:        "version",
			Description: "Prints server's version",
			Action:      version,
		},
		{
			Name:        "server",
			Description: "Start weaver server",
			Action:      startWeaver,
		},
	}

	app.Run(os.Args)
}
