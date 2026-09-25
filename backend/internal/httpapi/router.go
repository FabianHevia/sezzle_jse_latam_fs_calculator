package httpapi

import (
	"net/http"
)

// NewRouter constructs and configures the HTTP router with routes and middleware.
func NewRouter(allowedOrigin string) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/health", HandleHealth)
	mux.HandleFunc("/api/add", HandleAdd)
	mux.HandleFunc("/api/subtract", HandleSubtract)
	mux.HandleFunc("/api/multiply", HandleMultiply)
	mux.HandleFunc("/api/divide", HandleDivide)
	mux.HandleFunc("/api/power", HandlePower)
	mux.HandleFunc("/api/sqrt", HandleSqrt)
	mux.HandleFunc("/api/percentage", HandlePercentage)

	mux.HandleFunc("/", HandleNotFound)

	var handler http.Handler = mux
	handler = CORSMiddleware(allowedOrigin)(handler)
	handler = RecoverMiddleware(handler)
	handler = LoggingMiddleware(handler)

	return handler
}
