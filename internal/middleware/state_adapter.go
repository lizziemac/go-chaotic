package middleware

import (
	"net/http"
	"time"
)

// RequestState holds per-request state that can be used by handlers and middleware.
// It is initialized by the InitState middleware and can be accessed using GetState.
// Add fields as needed to hold request-specific data.
type RequestState struct {
	Start time.Time // Time when the request was started
	// Body  []byte
	//...
}

type StateHandler func(http.ResponseWriter, *http.Request, *RequestState)

// WithState ensures we can pass around state (RequestState) to all of our middleware
// without hiding the state in the request context. This makes it explicit that the handler
// needs the state, avoids the need to call GetState() inside the handler, and still enables
// compatibility with standard http.Handler for calls like http.HandleFunc or http.ListenAndServe.
func WithState(h StateHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		st := &RequestState{Start: time.Now()}
		h(w, r, st)
	})
}

// Adapt converts a standard http.Handler into a StateHandler so that we can keep using
// http.Handlers in our middleware chain while still passing the RequestState around.
func Adapt(h http.Handler) StateHandler {
	return func(w http.ResponseWriter, r *http.Request, st *RequestState) {
		h.ServeHTTP(w, r)
	}
}
