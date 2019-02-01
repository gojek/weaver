package instrumentation

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gojekfarm/weaver/internal/config"
	"github.com/gojekfarm/weaver/pkg/logger"
	statsd "gopkg.in/alexcesaro/statsd.v2"
)

var statsD *statsd.Client

func InitiateStatsDMetrics() error {
	statsDConfig := config.StatsD()

	if statsDConfig.Enabled() {
		flushPeriod := time.Duration(statsDConfig.FlushPeriodInSeconds()) * time.Second
		address := fmt.Sprintf("%s:%d", statsDConfig.Host(), statsDConfig.Port())

		var err error
		statsD, err = statsd.New(statsd.Address(address),
			statsd.Prefix(statsDConfig.Prefix()), statsd.FlushPeriod(flushPeriod))

		if err != nil {
			logger.Errorf("StatsD: Error initiating client %s", err)
			return err
		}

		logger.Infof("StatsD: Sending metrics")
	}

	return nil
}

func StatsDClient() *statsd.Client {
	return statsD
}

func CloseStatsDClient() {
	if statsD != nil {
		logger.Infof("StatsD: Shutting down")
		statsD.Close()
	}
}

func NewTiming() statsd.Timing {
	if statsD != nil {
		return statsD.NewTiming()
	}

	return statsd.Timing{}
}

func IncrementTotalRequestCount() {
	incrementProbe("request.total.count")
}

func IncrementAPIRequestCount(apiName string) {
	incrementProbe(fmt.Sprintf("request.api.%s.count", apiName))
}

func IncrementAPIStatusCount(apiName string, httpStatusCode int) {
	incrementProbe(fmt.Sprintf("request.api.%s.status.%d.count", apiName, httpStatusCode))
}

func IncrementAPIBackendRequestCount(apiName, backendName string) {
	incrementProbe(fmt.Sprintf("request.api.%s.backend.%s.count", apiName, backendName))
}

func IncrementAPIBackendStatusCount(apiName, backendName string, httpStatusCode int) {
	incrementProbe(fmt.Sprintf("request.api.%s.backend.%s.status.%d.count", apiName, backendName, httpStatusCode))
}

func IncrementCrashCount() {
	incrementProbe("request.internal.crash.count")
}

func IncrementNotFound() {
	incrementProbe(fmt.Sprintf("request.internal.%d.count", http.StatusNotFound))
}

func IncrementInternalAPIStatusCount(aclName string, statusCode int) {
	incrementProbe(fmt.Sprintf("request.api.%s.internal.status.%d.count", aclName, statusCode))
}

func TimeTotalLatency(timing statsd.Timing) {
	if statsD != nil {
		timing.Send("request.time.total")
	}

	return
}

func TimeAPILatency(apiName string, timing statsd.Timing) {
	if statsD != nil {
		timing.Send(fmt.Sprintf("request.api.%s.time.total", apiName))
	}

	return
}

func TimeAPIBackendLatency(apiName, backendName string, timing statsd.Timing) {
	if statsD != nil {
		timing.Send(fmt.Sprintf("request.api.%s.backend.%s.time.total", apiName, backendName))
	}

	return
}

func incrementProbe(key string) {
	if statsD == nil {
		return
	}

	go statsD.Increment(key)
}
