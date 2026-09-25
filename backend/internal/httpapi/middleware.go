package httpapi

import (
	"log"
	"net/http"
	"os"
	"time"
)

// CORSMiddleware manages Access-Control headers for cross-origin requests.
func CORSMiddleware(allowedOrigin string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RecoverMiddleware handles panics gracefully without crashing the server process.
func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("[PANIC RECOVERY] %v", err)
				writeError(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "an unexpected server error occurred")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

type responseLogger struct {
	http.ResponseWriter
	statusCode int
}

func (rl *responseLogger) WriteHeader(code int) {
	rl.statusCode = code
	rl.ResponseWriter.WriteHeader(code)
}

// LoggingMiddleware logs incoming HTTP requests to standard output.
func LoggingMiddleware(next http.Handler) http.Handler {
	logger := log.New(os.Stdout, "[HTTP] ", log.LstdFlags)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rl := &responseLogger{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rl, r)

		logger.Printf("%s %s %d %v", r.Method, r.URL.Path, rl.statusCode, time.Since(start))
	})
}
