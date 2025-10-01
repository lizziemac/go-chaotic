package middleware

import (
	"bytes"
	"io"
	"net/http"

	"example.com/myapp/internal/logger"
	"example.com/myapp/internal/middleware"
)

// LogRequestBody is a middleware that logs the HTTP request body along with the method and path.
// It reads the request body, logs its contents using the application's logger, and then resets
// the body so that downstream handlers can still access it.
//
// TODO: Move to internal/middleware (common) when truly specific API middleware is added. This is for demonstration.
// TODO: Enforce a size limit for production.
func LogRequestBody(next middleware.StateHandler) middleware.StateHandler {
	return func(w http.ResponseWriter, r *http.Request, st *middleware.RequestState) {
		if r.Body != nil {
			body, _ := io.ReadAll(r.Body)
			logger.Info("request received",
				"method", r.Method,
				"path", r.URL.Path,
				"body", body,
			)
			// reset the body, otherwise it would be empty and unusable by other parts of the server
			r.Body = io.NopCloser(bytes.NewBuffer(body))
		}
		// passes the request down the chain
		next(w, r, st)
	}
}
