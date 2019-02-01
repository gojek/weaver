package server

import (
	"net/http"

	"github.com/gojekfarm/weaver/internal/config"
	"github.com/gojekfarm/weaver/internal/httperror"
	"github.com/gojekfarm/weaver/pkg/instrumentation"
	"github.com/gojekfarm/weaver/pkg/logger"
	newrelic "github.com/newrelic/go-agent"
)

type proxy struct {
	router *Router
}

func (proxy *proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rw := &wrapperResponseWriter{ResponseWriter: w}

	if r.URL.Path == "/ping" || r.URL.Path == "/" {
		proxy.pingHandler(rw, r)
		return
	}

	timing := instrumentation.NewTiming()

	defer instrumentation.TimeTotalLatency(timing)
	instrumentation.IncrementTotalRequestCount()

	acl, err := proxy.router.Route(r)
	if err != nil || acl == nil {
		logger.Errorrf(r, "failed to find route: %+v for request: %s", err, r.URL.String())

		httperror.NotFoundHandler(rw, r)
		return
	}

	backend, err := acl.Endpoint.Shard(r)
	if backend == nil || err != nil {
		logger.Errorrf(r, "failed to find backend for acl %s for: %s, error: %s", acl.ID, r.URL.String(), err)

		httperror.Err503Handler{ACLName: acl.ID}.ServeHTTP(rw, r)
		return
	}

	instrumentation.IncrementAPIBackendRequestCount(acl.ID, backend.Name)

	instrumentation.IncrementAPIRequestCount(acl.ID)
	apiTiming := instrumentation.NewTiming()
	defer instrumentation.TimeAPILatency(acl.ID, apiTiming)

	apiBackendTiming := instrumentation.NewTiming()
	defer instrumentation.TimeAPIBackendLatency(acl.ID, backend.Name, apiBackendTiming)

	var s newrelic.ExternalSegment
	if txn, ok := w.(newrelic.Transaction); ok {
		s = newrelic.StartExternalSegment(txn, r)
	}
	backend.Handler.ServeHTTP(rw, r)

	s.End()

	logger.ProxyInfo(acl.ID, backend.Server.String(), r, rw.statusCode, rw)
	instrumentation.IncrementAPIStatusCount(acl.ID, rw.statusCode)
	instrumentation.IncrementAPIBackendStatusCount(acl.ID, backend.Name, rw.statusCode)
}

func (proxy *proxy) pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

func wrapNewRelicHandler(proxy *proxy) http.Handler {
	if !config.NewRelicConfig().Enabled {
		return proxy
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if path == "/ping" {
			proxy.ServeHTTP(w, r)
			return
		}

		_, next := newrelic.WrapHandleFunc(instrumentation.NewRelicApp(), path,
			func(w http.ResponseWriter, r *http.Request) {
				proxy.ServeHTTP(w, r)
			})

		next(w, r)
	})
}
