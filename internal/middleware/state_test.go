package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"example.com/myapp/internal/middleware"
)

func TestWithState_PassesState(t *testing.T) {
	var called bool        // global to verify handler was called
	var gotStart time.Time // global to capture the Start time from RequestState

	// Create a handler that checks the RequestState
	h := middleware.WithState(func(w http.ResponseWriter, r *http.Request, st *middleware.RequestState) {
		called = true
		gotStart = st.Start
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	if !called {
		t.Fatal("expected handler to be called")
	}
	if gotStart.IsZero() {
		t.Error("expected RequestState.Start to be set")
	}
	if status := w.Result().StatusCode; status != http.StatusOK {
		t.Errorf("expected status 200, got %d", status)
	}
}

func TestAdapt_WrapsHttpHandler(t *testing.T) {
	var called bool // global to verify handler was called

	// Create a standard http.Handler that will be adapted to StateHandler
	stdHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusNoContent)
	})

	// Adapt the standard http.Handler to a StateHandler
	h := middleware.Adapt(stdHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	// adapted is a StateHandler, call it with a dummy state
	h(w, req, &middleware.RequestState{Start: time.Now()})

	if !called {
		t.Fatal("expected adapted http.Handler to be called")
	}
	if status := w.Result().StatusCode; status != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", status)
	}
}
