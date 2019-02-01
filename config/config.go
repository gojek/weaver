package config

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	etcd "github.com/coreos/etcd/client"
	newrelic "github.com/newrelic/go-agent"
	"github.com/spf13/viper"
)

var appConfig Config

type Config struct {
	proxyHost       string
	proxyPort       int
	etcdKeyPrefix   string
	loggerLevel     string
	etcdEndpoints   []string
	etcdDialTimeout time.Duration
	statsDConfig    StatsDConfig
	newRelicConfig  newrelic.Config
	sentryDSN       string

	serverReadTimeout  time.Duration
	serverWriteTimeout time.Duration

	proxyConfig ProxyConfig
}

func Load() {
	viper.SetDefault("LOGGER_LEVEL", "error")
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("PROXY_PORT", "8081")

	viper.SetConfigName("weaver.conf")

	viper.AddConfigPath("./")
	viper.AddConfigPath("../")
	viper.AddConfigPath("../../")
	viper.SetConfigType("yaml")

	viper.ReadInConfig()
	viper.AutomaticEnv()

	appConfig = Config{
		proxyHost:          extractStringValue("PROXY_HOST"),
		proxyPort:          extractIntValue("PROXY_PORT"),
		etcdKeyPrefix:      extractStringValue("ETCD_KEY_PREFIX"),
		loggerLevel:        extractStringValue("LOGGER_LEVEL"),
		etcdEndpoints:      strings.Split(extractStringValue("ETCD_ENDPOINTS"), ","),
		etcdDialTimeout:    time.Duration(extractIntValue("ETCD_DIAL_TIMEOUT")),
		statsDConfig:       loadStatsDConfig(),
		newRelicConfig:     loadNewRelicConfig(),
		proxyConfig:        loadProxyConfig(),
		sentryDSN:          extractStringValue("SENTRY_DSN"),
		serverReadTimeout:  time.Duration(extractIntValue("SERVER_READ_TIMEOUT")),
		serverWriteTimeout: time.Duration(extractIntValue("SERVER_WRITE_TIMEOUT")),
	}
}

func ServerReadTimeoutInMillis() time.Duration {
	return appConfig.serverReadTimeout * time.Millisecond
}

func ServerWriteTimeoutInMillis() time.Duration {
	return appConfig.serverWriteTimeout * time.Millisecond
}

func ProxyServerAddress() string {
	return fmt.Sprintf("%s:%d", appConfig.proxyHost, appConfig.proxyPort)
}

func ETCDKeyPrefix() string {
	return appConfig.etcdKeyPrefix
}

func NewRelicConfig() newrelic.Config {
	return appConfig.newRelicConfig
}

func SentryDSN() string {
	return appConfig.sentryDSN
}

func StatsD() StatsDConfig {
	return appConfig.statsDConfig
}

func Proxy() ProxyConfig {
	return appConfig.proxyConfig
}

func NewETCDClient() (etcd.Client, error) {
	return etcd.New(etcd.Config{
		Endpoints:               appConfig.etcdEndpoints,
		HeaderTimeoutPerRequest: appConfig.etcdDialTimeout * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 10 * time.Second,
		},
	})
}

func LogLevel() string {
	return appConfig.loggerLevel
}

func extractStringValue(key string) string {
	checkPresenceOf(key)
	return viper.GetString(key)
}

func extractBoolValue(key string) bool {
	checkPresenceOf(key)
	return viper.GetBool(key)
}

func extractBoolValueDefaultToFalse(key string) bool {
	if !viper.IsSet(key) {
		return false
	}

	return viper.GetBool(key)
}

func extractIntValue(key string) int {
	checkPresenceOf(key)
	v, err := strconv.Atoi(viper.GetString(key))
	if err != nil {
		panic(fmt.Sprintf("key %s is not a valid Integer value", key))
	}

	return v
}

func checkPresenceOf(key string) {
	if !viper.IsSet(key) {
		panic(fmt.Sprintf("key %s is not set", key))
	}
}
