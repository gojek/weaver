package httperror

import (
	"encoding/json"
	"net/http"

	"github.com/gojektech/weaver/pkg/instrumentation"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	instrumentation.IncrementNotFound()
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusNotFound)

	errorResponse := WeaverResponse{
		Errors: []ErrorDetails{
			ErrorDetails{
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
