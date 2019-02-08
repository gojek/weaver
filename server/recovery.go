package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	raven "github.com/getsentry/raven-go"
	"github.com/gojektech/weaver/pkg/instrumentation"
	"github.com/gojektech/weaver/pkg/logger"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		defer func() {
			if err := recover(); err != nil {
				instrumentation.IncrementCrashCount()

				var recoveredErr error
				switch val := err.(type) {
				case error:
					recoveredErr = val
				case string:
					recoveredErr = fmt.Errorf(val)
				}

				raven.CaptureError(recoveredErr, map[string]string{"error": recoveredErr.Error(), "request_url": r.URL.String()})

				//TODO: move this as a method inside error.go
				logger.Errorrf(r, "failed to route request: %+v", err)
				errorResponse := weaverResponse{
					Errors: []errorDetails{
						{
							Code:            "weaver:service:unavailable",
							Message:         "Something went wrong",
							MessageTitle:    "Failure",
							MessageSeverity: "failure",
						},
					},
				}

				response, _ := json.Marshal(errorResponse)
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write(response)
				return
			}
		}()

		next.ServeHTTP(w, r)
	})
}
