package handlers

import (
	"encoding/json"
	"net/http"
)

// ContextResponseWriter implements the ResponseWriter interface for the context package
type ContextResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (crw *ContextResponseWriter) init() {
	crw.statusCode = http.StatusOK
}

// WriteHeader sets the header for the ResponseWriter interface in the context package
func (crw *ContextResponseWriter) WriteHeader(code int) {
	crw.statusCode = code
	crw.ResponseWriter.WriteHeader(code)
}

// GetStatus sets the status code for the ResponseWriter interface in the context package
func (crw *ContextResponseWriter) GetStatus() int {
	return crw.statusCode
}

// Response for all JSON requests
type Response struct {
	context ContextResponseWriter
}

// New instantiates a new Response struct and attaches the Gin context.
// It returns the Response struct.
func New(c http.ResponseWriter) *Response {
	r := new(Response)
	r.context = ContextResponseWriter{statusCode: http.StatusOK, ResponseWriter: c}
	r.context.Header().Set("Content-Type", "application/json; charset=utf-8")
	return r
}

// Render encodes the status and content to the Response struct
func (r *Response) Render(code int, content interface{}) {
	r.context.WriteHeader(code)
	json.NewEncoder(r.context.ResponseWriter).Encode(content)
}
