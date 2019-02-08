package server

import (
	"encoding/json"
	"net/http"

	"github.com/gojektech/weaver/pkg/instrumentation"
)

type weaverResponse struct {
	Errors []errorDetails `json:"errors"`
}

type errorDetails struct {
	Code            string `json:"code"`
	Message         string `json:"message"`
	MessageTitle    string `json:"message_title"`
	MessageSeverity string `json:"message_severity"`
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	instrumentation.IncrementNotFound()
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusNotFound)

	errorResponse := weaverResponse{
		Errors: []errorDetails{
			{
				Code:            "weaver:route:not_found",
				Message:         "Something went wrong",
				MessageTitle:    "Failure",
				MessageSeverity: "failure",
			},
		},
	}

	response, _ := json.Marshal(errorResponse)
	w.Write(response)
}

func internalServerError(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusInternalServerError)

	errorResponse := weaverResponse{
		Errors: []errorDetails{
			{
				Code:            "weaver:service:unavailable",
				Message:         "Something went wrong",
				MessageTitle:    "Internal error",
				MessageSeverity: "failure",
			},
		},
	}

	response, _ := json.Marshal(errorResponse)
	w.Write(response)
}

// TODO: decouple instrumentation from this errors function
type err503Handler struct {
	ACLName string
}

func (eh err503Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	failureHTTPStatus := http.StatusServiceUnavailable
	instrumentation.IncrementInternalAPIStatusCount(eh.ACLName, failureHTTPStatus)

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
	w.WriteHeader(failureHTTPStatus)
	w.Write(response)
	return
}
