// Package server defines and manages the HTTP server for the API.
//
// It provides setup and configuration for the HTTP server, including:
//   - Route registration
//   - Middleware chains
//   - Graceful shutdown handling
//
// The server package acts as the entry point for serving HTTP traffic
// and ties together handlers, middleware, and configuration.
//
// Example usage:
//
//	func main() {
//	    srv := &http.Server{
//	        Addr:    ":8080",
//	        Handler: server.Setup(),
//	    }
//
//	    log.Printf("listening on %s", srv.Addr)
//	    if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
//	        log.Fatal(err)
//	    }
//	}
package server
