package server

import (
	"context"
	"log"
	"net/http"

	"github.com/gojekfarm/weaver/internal/config"
	"github.com/gojekfarm/weaver/internal/middleware"
	"github.com/gojekfarm/weaver/pkg/util"
)

var server *Weaver

type Weaver struct {
	httpServer      *http.Server
	cancelRouteSync context.CancelFunc
}

func ShutdownServer(ctx context.Context) {
	server.cancelRouteSync()
	server.httpServer.Shutdown(ctx)
}

func StartServer() {
	routeSyncCtx, cancelRouteSync := context.WithCancel(context.Background())
	routeLoader, err := NewETCDRouteLoader()
	if err != nil {
		log.Printf("StartServer: failed to initialise etcd route loader: %s", err)
		cancelRouteSync()
	}

	proxyRouter := NewRouter(routeLoader)
	err = proxyRouter.BootstrapRoutes(context.Background())
	if err != nil {
		log.Printf("StartServer: failed to initialise proxy router: %s", err)
		cancelRouteSync()
	}

	log.Printf("StartServer: bootstraped routes from etcd")

	go proxyRouter.WatchRouteUpdates(routeSyncCtx)

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
		httpServer:      httpServer,
		cancelRouteSync: cancelRouteSync,
	}

	log.Printf("StartServer: starting weaver on %s", server.httpServer.Addr)
	log.Printf("Keep-Alive: %s", util.BoolToOnOff(keepAliveEnabled))

	if err := server.httpServer.ListenAndServe(); err != nil {
		log.Fatalf("StartServer: starting weaver failed with %s", err)
	}
}
