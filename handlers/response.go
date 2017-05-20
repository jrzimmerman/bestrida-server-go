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
	context    ContextResponseWriter
	Success    bool        `json:"success"`
	StatusCode int         `json:"status_code"`
	StatusText string      `json:"status_text"`
	Content    interface{} `json:"content"`
}

// New instantiates a new Response struct and attaches the Gin context.
// It returns the Response struct.
func New(c http.ResponseWriter) *Response {
	r := new(Response)
	r.context = ContextResponseWriter{statusCode: http.StatusOK, ResponseWriter: c}
	r.context.Header().Set("Content-Type", "application/json; charset=utf-8")
	r.Success = false
	r.StatusCode = http.StatusInternalServerError
	return r
}

// SetResponse sets the response status code and content.
func (r *Response) SetResponse(code int, content interface{}) {
	r.StatusCode = code
	r.Content = content
}

// Render encodes the status and content to the Response struct
func (r *Response) Render() {
	if r.StatusCode >= 200 && r.StatusCode < 300 {
		r.Success = true
	}

	r.StatusText = http.StatusText(r.StatusCode)
	r.context.WriteHeader(r.StatusCode)
	json.NewEncoder(r.context.ResponseWriter).Encode(r)
}
