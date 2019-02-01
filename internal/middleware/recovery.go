package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"

	raven "github.com/getsentry/raven-go"
	"github.com/gojekfarm/weaver/internal/httperror"
	"github.com/gojekfarm/weaver/pkg/instrumentation"
	"github.com/gojekfarm/weaver/pkg/logger"
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

				logger.Errorrf(r, "failed to route request: %+v", err)
				errorResponse := httperror.WeaverResponse{
					Errors: []httperror.ErrorDetails{
						httperror.ErrorDetails{
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
