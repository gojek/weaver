package config

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShouldLoadConfigFromFile(t *testing.T) {
	Load()

	assert.NotEmpty(t, LogLevel())
	assert.NotNil(t, loadStatsDConfig().Prefix())
	assert.NotNil(t, loadStatsDConfig().FlushPeriodInSeconds())
	assert.NotNil(t, loadStatsDConfig().Port())
	assert.NotNil(t, loadStatsDConfig().Enabled())
}

func TestShouldLoadFromEnvVars(t *testing.T) {
	configVars := map[string]string{
		"LOGGER_LEVEL":                   "info",
		"NEW_RELIC_APP_NAME":             "newrelic",
		"NEW_RELIC_LICENSE_KEY":          "licence",
		"NEW_RELIC_ENABLED":              "true",
		"STATSD_PREFIX":                  "weaver",
		"STATSD_FLUSH_PERIOD_IN_SECONDS": "20",
		"STATSD_HOST":                    "statsd",
		"STATSD_PORT":                    "8125",
		"STATSD_ENABLED":                 "true",
		"ETCD_KEY_PREFIX":                "weaver",

		"PROXY_DIALER_TIMEOUT_IN_MS":    "10",
		"PROXY_DIALER_KEEP_ALIVE_IN_MS": "10",
		"PROXY_MAX_IDLE_CONNS":          "200",
		"PROXY_IDLE_CONN_TIMEOUT_IN_MS": "20",
		"SENTRY_DSN":                    "dsn",

		"SERVER_READ_TIMEOUT":  "100",
		"SERVER_WRITE_TIMEOUT": "100",
	}

	for k, v := range configVars {
		err := os.Setenv(k, v)
		require.NoError(t, err, fmt.Sprintf("failed to set env for %s key", k))
	}

	Load()

	expectedStatsDConfig := StatsDConfig{
		prefix:               "weaver",
		flushPeriodInSeconds: 20,
		host:                 "statsd",
		port:                 8125,
		enabled:              true,
	}

	assert.Equal(t, "info", LogLevel())

	assert.Equal(t, "newrelic", loadNewRelicConfig().AppName)
	assert.Equal(t, "licence", loadNewRelicConfig().License)
	assert.True(t, loadNewRelicConfig().Enabled)

	assert.Equal(t, expectedStatsDConfig, loadStatsDConfig())
	assert.Equal(t, "weaver", ETCDKeyPrefix())
	assert.Equal(t, "dsn", SentryDSN())

	assert.Equal(t, time.Duration(10)*time.Millisecond, Proxy().ProxyDialerTimeoutInMS())
	assert.Equal(t, time.Duration(10)*time.Millisecond, Proxy().ProxyDialerKeepAliveInMS())
	assert.Equal(t, 200, Proxy().ProxyMaxIdleConns())
	assert.Equal(t, time.Duration(20)*time.Millisecond, Proxy().ProxyIdleConnTimeoutInMS())

	assert.Equal(t, time.Duration(100)*time.Millisecond, ServerReadTimeoutInMillis())
	assert.Equal(t, time.Duration(100)*time.Millisecond, ServerWriteTimeoutInMillis())
}
