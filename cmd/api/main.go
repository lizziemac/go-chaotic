package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"example.com/myapp/internal/logger"
	"example.com/myapp/internal/middleware"
	commonmw "example.com/myapp/internal/middleware"
	reverseproxy "example.com/myapp/internal/proxy"
	"example.com/myapp/internal/services"
)

func main() {
	// Create a context that closes when the app is interrupted
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	var config services.ConfigStore = services.NewConfigRegistry(ctx)

	proxy := reverseproxy.Setup(ctx, config) // set up the reverse proxy

	proxyStateHandler := func(w http.ResponseWriter, r *http.Request, st *middleware.RequestState) {
		proxy.ServeHTTP(w, r)
	}

	h := commonmw.LogRequestBody(proxyStateHandler)
	h = commonmw.RequestTimer(h)
	// add more middleware here
	// e.g. `h = middleware.X(h)`

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
