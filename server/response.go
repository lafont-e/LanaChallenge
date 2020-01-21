package server

import (
	"encoding/json"
	"net/http"
)

type Data struct {
	Type    string
	Content interface{}
}

// Response - A standardised response format
type Response struct {
	Status  string `json:"status"`         // Can be 'ok' or 'fail'
	Code    int    `json:"code"`           // Any valid HTTP response code
	Message string `json:"message"`        // Any relevant message (optional)
	Data    *Data  `json:"data,omitempty"` // Data to pass along to the response (optional)
}

// New returns a new Response
func NewResponse(code int, message string, data *Data) *Response {
	var status string = StatusFail
	if code >= http.StatusOK && code < http.StatusBadRequest {
		status = StatusOk
	}

	return &Response{
		Code:    code,
		Status:  status,
		Message: message,
		Data:    data,
	}
}

// WriteTo - response writer to write the default json response to.
func (r *Response) WriteTo(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.Code)

	// Don't write a body for 204s.
	if r.Code == http.StatusNoContent {
		return nil
	}

	j, err := json.Marshal(r)
	if err != nil {
		return err
	}
	_, err = w.Write(j)
	return err
}
