package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"example.com/myapp/internal/api/handlers"
	"example.com/myapp/internal/logger"
	"example.com/myapp/internal/middleware"
	reverseproxy "example.com/myapp/internal/proxy"
	"example.com/myapp/internal/services"
)

func main() {
	// Create a context that closes when the app is interrupted
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	var config services.ConfigStore = services.NewConfigRegistry(ctx)

	mux := http.NewServeMux()

	// Create a router for reverse-proxy specific APIs
	mux.HandleFunc("PUT /proxy/api/v1/config", handlers.UpsertConfig(config))
	mux.HandleFunc("GET /proxy/api/v1/config", handlers.GetConfig(config))

	// Catch-all for /proxy/api/v1/config to prevent fall-through to the proxy
	mux.HandleFunc("/proxy/api/v1/config", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	})

	// Set up the reverse proxy
	proxy := reverseproxy.Setup(ctx, config)
	mux.Handle("/", proxy) // route everything else via the reverse proxy

	// Setup the middleware to apply to all routes
	h := middleware.Adapt(mux)
	h = middleware.LogRequestBody(h)
	h = middleware.RequestTimer(h)
	// ...

	// Finalize middleware by applying a state context
	finalHandler := middleware.WithState(h)

	srv := &http.Server{
		Addr:    ":8081",
		Handler: finalHandler,
	}

	// Start the server in a separate goroutine
	go func() {
		logger.Info("chaos proxy starting", "port", 8081)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("failed to start chaos proxy", "error", err)
			cancel() // Tell the main thread to shut down immediately
		}
	}()

	// Wait for the context to be canceled (e.g. by Ctrl+C)
	<-ctx.Done()
	logger.Info("shutting down...")

	// Create a timeout context to force a shutdown if requests take too long to finish
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Fatal("server forced to shutdown", "error", err)
	}
}
