package server

import (
	"encoding/json"
	"net/http"

	"github.com/gojektech/weaver/pkg/instrumentation"
)

type Err503Handler struct {
	ACLName string
}

func (eh Err503Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	failureHTTPStatus := http.StatusServiceUnavailable
	instrumentation.IncrementInternalAPIStatusCount(eh.ACLName, failureHTTPStatus)

	errorResponse := WeaverResponse{
		Errors: []ErrorDetails{
			ErrorDetails{
				Code:            "weaver:service:unavailable",
				Message:         "Something went wrong",
				MessageTitle:    "Failure",
				MessageSeverity: "failure",
			},
		},
	}

	response, _ := json.Marshal(errorResponse)
	w.WriteHeader(failureHTTPStatus)
	w.Write(response)
	return
}
