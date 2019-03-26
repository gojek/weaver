package server

import (
	"context"
	"github.com/gojektech/weaver/internal/cli"
	"github.com/gojektech/weaver/pkg/logger"
	"github.com/gojektech/weaver/server"
	"os"
	"os/signal"
	"syscall"
)

const (
	startCmdName        = "start"
	startCmdUsage       = "Run Weaver server"
	startCmdDescription = "Run Weaver server"
)

var serverStartCmd = cli.NewDefaultCommand(startCmdName, startCmdUsage, startCmdDescription, startServer)

func startServer(c *cli.Context) error {
	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	go server.StartServer(ctx, c.RouteLoader)

	sig := <-sigC
	logger.Infof("Received %d, shutting down", sig)

	defer cancel()
	server.ShutdownServer(ctx)

	return nil
}

func init() {
	weaverServerCmd.RegisterCommand(serverStartCmd)
}
