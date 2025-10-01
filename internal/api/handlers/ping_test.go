package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/myapp/internal/services"
)

func TestPingHttpHandling(t *testing.T) {
	// Generate w and r
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()

	// Call the Ping handler
	Ping(w, req)

	// Check the status code
	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", w.Result().Status)
	}

	// Check the response body
	var resp PingResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}
	expectedMessage := services.Pong()
	if resp.Message != expectedMessage {
		t.Errorf("expected message %q; got %q", expectedMessage, resp.Message)
	}
}
