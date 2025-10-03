package server

import (
	"net/http"

	"example.com/myapp/internal/api/handlers"
	apimw "example.com/myapp/internal/api/middleware"
	commonmw "example.com/myapp/internal/middleware"
)

func Setup() http.Handler {
	// API Routes
	apiMux := http.NewServeMux()                         // create a new mux for API routes (so we can apply API-specific middleware)
	apiMux.HandleFunc("GET /api/v1/ping", handlers.Ping) // register API handler for /api/v1/ping

	// API Middleware
	h := apimw.LogRequestBody(commonmw.Adapt(apiMux)) // adapt apiMux to StateHandler and wrap with LogRequestBody

	// Root Mux
	rootMux := http.NewServeMux()                  // create a root mux for all routes (this will be where common middleware is applied)
	rootMux.Handle("/api/", commonmw.WithState(h)) // register the API mux under /api/ with state middleware

	// Common Middleware
	h = commonmw.Adapt(rootMux)  // adapt rootMux to StateHandler so we can wrap with common middleware
	h = commonmw.RequestTimer(h) // wrap with RequestTimer middleware

	return commonmw.WithState(h)
}
