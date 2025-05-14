package middleware 

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type responseWriter struct {
	http.ResponseWriter
	status int
	size int	
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.status == 0 {
		rw.status = 200
	}
	n, err := rw.ResponseWriter.Write(b)
	rw.size += n

	return n, err
}

func ZapLogger(logger *zap.SugaredLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := &responseWriter{ResponseWriter: w}
			reqID := middleware.GetReqID(r.Context())
			
			start := time.Now()
			next.ServeHTTP(rw, r)
			elapsed := time.Since(start)

			if rw.status == 0 {
				rw.status = http.StatusOK
			}

			logger.Infow("HTTP request",
				"request_id", reqID,
				"method", r.Method,
				"path", r.URL.Path,
				"status", rw.status,
				"bytes", rw.size,
				"duration", elapsed.String(),
				"remote_addr", r.RemoteAddr,
				"user_agent", r.UserAgent(),
			)
		})
	}
}