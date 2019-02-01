package config

type StatsDConfig struct {
	prefix               string
	flushPeriodInSeconds int
	host                 string
	port                 int
	enabled              bool
}

func loadStatsDConfig() StatsDConfig {
	return StatsDConfig{
		prefix:               extractStringValue("STATSD_PREFIX"),
		flushPeriodInSeconds: extractIntValue("STATSD_FLUSH_PERIOD_IN_SECONDS"),
		host:                 extractStringValue("STATSD_HOST"),
		port:                 extractIntValue("STATSD_PORT"),
		enabled:              extractBoolValue("STATSD_ENABLED"),
	}
}

func (sdc StatsDConfig) Prefix() string {
	return sdc.prefix
}

func (sdc StatsDConfig) FlushPeriodInSeconds() int {
	return sdc.flushPeriodInSeconds
}

func (sdc StatsDConfig) Host() string {
	return sdc.host
}

func (sdc StatsDConfig) Port() int {
	return sdc.port
}

func (sdc StatsDConfig) Enabled() bool {
	return sdc.enabled
}
