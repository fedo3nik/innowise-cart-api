package models

import (
	"fmt"
	"net/http"
)

type ResponseError struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

// in case of logging or calling this error
func (e ResponseError) Error() string {
	return fmt.Sprintf("Err message: %s, statusCode: %d", e.Message, e.StatusCode)
}

func (e ResponseError) ShowError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(e.StatusCode)
	fmt.Fprintln(w, e.Message)
}

func NewResponseError(code int, msg string) *ResponseError {
	return &ResponseError{
		StatusCode: code,
		Message:    msg,
	}

}
