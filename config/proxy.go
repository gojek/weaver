package config

import "time"

type ProxyConfig struct {
	proxyDialerTimeoutInMS   int
	proxyDialerKeepAliveInMS int
	proxyMaxIdleConns        int
	proxyIdleConnTimeoutInMS int
	keepAliveEnabled         bool
}

func loadProxyConfig() ProxyConfig {
	return ProxyConfig{
		proxyDialerTimeoutInMS:   extractIntValue("PROXY_DIALER_TIMEOUT_IN_MS"),
		proxyDialerKeepAliveInMS: extractIntValue("PROXY_DIALER_KEEP_ALIVE_IN_MS"),
		proxyMaxIdleConns:        extractIntValue("PROXY_MAX_IDLE_CONNS"),
		proxyIdleConnTimeoutInMS: extractIntValue("PROXY_IDLE_CONN_TIMEOUT_IN_MS"),
		keepAliveEnabled:         extractBoolValueDefaultToFalse("PROXY_KEEP_ALIVE_ENABLED"),
	}
}

func (pc ProxyConfig) ProxyDialerTimeoutInMS() time.Duration {
	return time.Duration(pc.proxyDialerTimeoutInMS) * time.Millisecond
}

func (pc ProxyConfig) ProxyDialerKeepAliveInMS() time.Duration {
	return time.Duration(pc.proxyDialerKeepAliveInMS) * time.Millisecond
}

func (pc ProxyConfig) ProxyMaxIdleConns() int {
	return pc.proxyMaxIdleConns
}

func (pc ProxyConfig) ProxyIdleConnTimeoutInMS() time.Duration {
	return time.Duration(pc.proxyIdleConnTimeoutInMS) * time.Millisecond
}

func (pc ProxyConfig) KeepAliveEnabled() bool {
	return pc.keepAliveEnabled
}
