package middleware

import (
	"log"
	"net/http"
)

type LoggingMiddleware struct {
	handler http.Handler
}

func NewLoggingMiddleware(next http.Handler) *LoggingMiddleware {
	return &LoggingMiddleware{
		handler: next,
	}
}

func (l *LoggingMiddleware) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	l.handler.ServeHTTP(res, req)

	log.Printf("%s----------->%s", req.Method, req.URL.Path)
}
