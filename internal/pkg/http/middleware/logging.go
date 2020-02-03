package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func LoggingMiddleware(h http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bg := time.Now()
		logrus.Infof("path: %s, method: %s, request_time: %v", r.URL.Path, r.Method, bg)
		h.ServeHTTP(w, r)
		logrus.Infof("path: %s, response_time: %v", r.URL.Path, time.Since(bg))
	})
}
