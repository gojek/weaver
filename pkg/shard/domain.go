package shard

import (
	"fmt"
	"time"

	"github.com/gojektech/weaver"
	"github.com/gojektech/weaver/config"
	"github.com/pkg/errors"
)

type CustomError struct {
	ExitMessage string
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("[error]  %s", e.ExitMessage)
}

func Error(msg string) error {
	return &CustomError{msg}
}

type BackendDefinition struct {
	BackendName string   `json:"backend_name"`
	BackendURL  string   `json:"backend"`
	Timeout     *float64 `json:"timeout,omitempty"`
}

func (bd BackendDefinition) Validate() error {
	if bd.BackendName == "" {
		return errors.WithStack(fmt.Errorf("missing backend name in shard config: %+v", bd))
	}

	if bd.BackendURL == "" {
		return errors.WithStack(fmt.Errorf("missing backend url in shard config: %+v", bd))
	}

	return nil
}

func toBackends(shardConfig map[string]BackendDefinition) (map[string]*weaver.Backend, error) {
	backends := map[string]*weaver.Backend{}

	for key, backendDefinition := range shardConfig {
		if err := backendDefinition.Validate(); err != nil {
			return nil, errors.Wrapf(err, "failed to validate backend definition")
		}

		backend, err := parseBackend(backendDefinition)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parseBackends from backendDefinition")
		}

		backends[key] = backend
	}

	return backends, nil
}

func parseBackend(shardConfig BackendDefinition) (*weaver.Backend, error) {
	timeoutInDuration := config.Proxy().ProxyDialerTimeoutInMS()

	if shardConfig.Timeout != nil {
		timeoutInDuration = time.Duration(*shardConfig.Timeout)
	}

	backendOptions := weaver.BackendOptions{
		Timeout: timeoutInDuration * time.Millisecond,
	}

	return weaver.NewBackend(shardConfig.BackendName, shardConfig.BackendURL, backendOptions)
}
