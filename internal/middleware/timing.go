package middleware

import (
	"net/http"
	"time"

	"example.com/myapp/internal/logger"
)

// RequestTimer is a middleware that logs the duration of each HTTP request.
// It captures the start time from the RequestState (initialized by WithState middleware adapter)
// and logs the elapsed time after the request has been processed.
func RequestTimer(next StateHandler) StateHandler {
	return func(w http.ResponseWriter, r *http.Request, st *RequestState) {
		next(w, r, st)
		if st != nil {
			elapsed := time.Since(st.Start)
			logger.Info("request complete",
				"method", r.Method,
				"path", r.URL.Path,
				"duration", elapsed,
			)
		} else {
			logger.Error("request complete but cannot record duration")
		}
	}
}
