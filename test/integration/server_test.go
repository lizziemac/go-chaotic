package integration

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/myapp/internal/api/server"
)

func TestNotFoundReturns404(t *testing.T) {
	srv := httptest.NewServer(server.Setup())
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/api/v1/does-not-exist")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected %d, got %d", http.StatusNotFound, resp.StatusCode)
	}
}

func TestPingDoesNotAcceptPOST(t *testing.T) {
	srv := httptest.NewServer(server.Setup())
	defer srv.Close()

	resp, err := http.Post(srv.URL+"/api/v1/ping", "application/json", bytes.NewBuffer([]byte(`{}`)))
	if err != nil {
		t.Fatalf("post failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected %d, got %d", http.StatusMethodNotAllowed, resp.StatusCode)
	}
	respBody, _ := io.ReadAll(resp.Body)
	if !bytes.Contains(respBody, []byte("method not allowed")) {
		t.Errorf("expected 'method not allowed' response, got %q", respBody)
	}
}
