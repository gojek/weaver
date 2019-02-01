package server

import "net/http"

type wrapperResponseWriter struct {
	statusCode int
	http.ResponseWriter
}

func (w *wrapperResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *wrapperResponseWriter) Write(data []byte) (int, error) {
	return w.ResponseWriter.Write(data)
}

func (w *wrapperResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
