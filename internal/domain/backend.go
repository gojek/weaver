package domain

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/gojekfarm/weaver/internal/config"
	"github.com/pkg/errors"
)

type Backend struct {
	Handler http.Handler
	Server  *url.URL
	Name    string
}

type BackendOptions struct {
	Timeout time.Duration
}

func NewBackend(name string, serverURL string, options BackendOptions) (*Backend, error) {
	server, err := url.Parse(serverURL)
	if err != nil {
		return nil, errors.Wrapf(err, "URL Parsing failed for: %s", serverURL)
	}

	return &Backend{
		Name:    name,
		Handler: newWeaverReverseProxy(server, options),
		Server:  server,
	}, nil
}

func newWeaverReverseProxy(target *url.URL, options BackendOptions) *httputil.ReverseProxy {
	proxyConfig := config.Proxy()

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   options.Timeout,
			KeepAlive: proxyConfig.ProxyDialerKeepAliveInMS(),
			DualStack: true,
		}).DialContext,

		MaxIdleConns:      proxyConfig.ProxyMaxIdleConns(),
		IdleConnTimeout:   proxyConfig.ProxyIdleConnTimeoutInMS(),
		DisableKeepAlives: !proxyConfig.KeepAliveEnabled(),
	}

	return proxy
}
