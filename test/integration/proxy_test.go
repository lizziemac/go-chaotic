package integration

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"example.com/myapp/internal/reverseproxy"
)

func TestProxyForwardsRequest(t *testing.T) {
	// Create a dummy backend server to receive the forwarded request
	backendCalled := false
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		backendCalled = true
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello from backend"))
	}))
	defer backend.Close()

	// Set the required env var for the reverse proxy setup
	os.Setenv("TARGET_URL", backend.URL)
	defer os.Unsetenv("TARGET_URL")

	// Set up the reverse proxy
	proxy := reverseproxy.Setup()

	// Create a test request hitting the proxy
	req := httptest.NewRequest(http.MethodGet, "/test-path", nil)
	w := httptest.NewRecorder()

	// Serve the request via the proxy
	proxy.ServeHTTP(w, req)

	// Assertions
	if !backendCalled {
		t.Fatal("expected the backend to be called, but it wasn't")
	}
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	expected := "hello from backend"
	if !bytes.Contains(w.Body.Bytes(), []byte(expected)) {
		t.Errorf("expected %q response, got %q", expected, w.Body.String())
	}
}
