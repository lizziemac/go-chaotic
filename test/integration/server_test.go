package integration

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/myapp/internal/api/server"
)

func TestNotFoundReturns404(t *testing.T) {
	handler := server.Setup()
	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected %d, got %d", http.StatusNotFound, w.Code)
	}
	expected := "404 page not found"
	if !bytes.Contains(w.Body.Bytes(), []byte(expected)) {
		t.Errorf("expected %q response, got %q", expected, w.Body.String())
	}
}

func TestPingDoesNotAcceptPOST(t *testing.T) {
	handler := server.Setup()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/ping", bytes.NewBuffer([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
	expected := http.StatusText(http.StatusMethodNotAllowed)
	if !bytes.Contains(w.Body.Bytes(), []byte(expected)) {
		t.Errorf("expected %q response, got %q", expected, w.Body.String())
	}
}
