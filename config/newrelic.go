package config

import (
	newrelic "github.com/newrelic/go-agent"
)

func loadNewRelicConfig() newrelic.Config {
	config := newrelic.NewConfig(extractStringValue("NEW_RELIC_APP_NAME"),
		extractStringValue("NEW_RELIC_LICENSE_KEY"))
	config.Enabled = extractBoolValue("NEW_RELIC_ENABLED")
	return config
}
