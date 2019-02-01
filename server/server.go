package server

import (
	"context"
	"log"
	"net/http"

	"github.com/gojektech/weaver/config"
	"github.com/gojektech/weaver/middleware"
	"github.com/gojektech/weaver/pkg/util"
)

var server *Weaver

type Weaver struct {
	httpServer *http.Server
}

func ShutdownServer(ctx context.Context) {
	server.httpServer.Shutdown(ctx)
}

func StartServer(ctx context.Context, routeLoader RouteLoader) {
	proxyRouter := NewRouter(routeLoader)
	err := proxyRouter.BootstrapRoutes(context.Background())
	if err != nil {
		log.Printf("StartServer: failed to initialise proxy router: %s", err)
	}

	log.Printf("StartServer: bootstraped routes from etcd")

	go proxyRouter.WatchRouteUpdates(ctx)

	proxy := middleware.Recover(wrapNewRelicHandler(&proxy{
		router: proxyRouter,
	}))

	httpServer := &http.Server{
		Addr:         config.ProxyServerAddress(),
		Handler:      proxy,
		ReadTimeout:  config.ServerReadTimeoutInMillis(),
		WriteTimeout: config.ServerWriteTimeoutInMillis(),
	}

	keepAliveEnabled := config.Proxy().KeepAliveEnabled()
	httpServer.SetKeepAlivesEnabled(keepAliveEnabled)

	server = &Weaver{
		httpServer: httpServer,
	}

	log.Printf("StartServer: starting weaver on %s", server.httpServer.Addr)
	log.Printf("Keep-Alive: %s", util.BoolToOnOff(keepAliveEnabled))

	if err := server.httpServer.ListenAndServe(); err != nil {
		log.Fatalf("StartServer: starting weaver failed with %s", err)
	}
}
