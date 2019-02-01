package instrumentation

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gojektech/weaver/config"
	newrelic "github.com/newrelic/go-agent"
)

type ctxKey int

const txKey ctxKey = 0

var newRelicApp newrelic.Application

func InitNewRelic() newrelic.Application {
	cfg := config.NewRelicConfig()
	if cfg.Enabled {
		app, err := newrelic.NewApplication(cfg)
		if err != nil {
			log.Fatalf(err.Error())
		}

		newRelicApp = app
	}
	return newRelicApp
}

func ShutdownNewRelic() {
	if config.NewRelicConfig().Enabled {
		newRelicApp.Shutdown(time.Second)
	}
}

func NewRelicApp() newrelic.Application {
	return newRelicApp
}

func StartRedisSegmentNow(op string, coll string, txn newrelic.Transaction) newrelic.DatastoreSegment {
	s := newrelic.DatastoreSegment{
		Product:    newrelic.DatastoreRedis,
		Collection: coll,
		Operation:  op,
	}

	s.StartTime = newrelic.StartSegmentNow(txn)
	return s
}

func NewContext(ctx context.Context, w http.ResponseWriter) context.Context {
	if config.NewRelicConfig().Enabled {
		tx, ok := w.(newrelic.Transaction)
		if !ok {
			return ctx
		}
		return context.WithValue(ctx, txKey, tx)
	}
	return ctx
}

func NewContextWithTransaction(ctx context.Context, tx newrelic.Transaction) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

func GetTx(ctx context.Context) (newrelic.Transaction, bool) {
	tx, ok := ctx.Value(txKey).(newrelic.Transaction)
	return tx, ok
}
