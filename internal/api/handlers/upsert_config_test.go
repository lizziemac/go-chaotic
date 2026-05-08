package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/myapp/internal/services"
)

// mockStore is a simple mock implementation of services.ConfigStore
type mockStore struct {
	upsertOk  bool
	getConfig *services.Config
}

func (m *mockStore) UpsertConfig(userID string, config services.Config) (*services.Config, bool) {
	if !m.upsertOk {
		return nil, false
	}
	return &config, true
}

func (m *mockStore) GetConfig(userID string) (*services.Config, bool) {
	if m.getConfig != nil {
		return m.getConfig, true
	}
	return nil, false
}

func TestUpsertConfig(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		body           string
		mockUpsertOk   bool
		expectedStatus int
	}{
		{
			name:           "successful upsert",
			userID:         "user-1",
			body:           `{"mode": 2, "latency_delay_ns": 5000000000}`,
			mockUpsertOk:   true,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "missing user id",
			userID:         "",
			body:           `{"mode": 2}`,
			mockUpsertOk:   true,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid json",
			userID:         "user-1",
			body:           `{"mode": "not-a-number"}`,
			mockUpsertOk:   true,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "store failure",
			userID:         "user-1",
			body:           `{"mode": 2}`,
			mockUpsertOk:   false,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &mockStore{upsertOk: tt.mockUpsertOk}
			handler := UpsertConfig(store)

			req := httptest.NewRequest(http.MethodPut, "/proxy/api/v1/config", bytes.NewBufferString(tt.body))
			if tt.userID != "" {
				req.Header.Set("X-User-ID", tt.userID)
			}
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				var resp services.Config
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response body: %v", err)
				}
				if resp.Mode != 2 {
					t.Errorf("expected mode 2, got %v", resp.Mode)
				}
			}
		})
	}
}

func TestUpsertConfigDoesNotAcceptPOST(t *testing.T) {
	store := &mockStore{upsertOk: true}

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /proxy/api/v1/config", UpsertConfig(store))

	req := httptest.NewRequest(http.MethodPost, "/proxy/api/v1/config", bytes.NewBuffer([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
	expected := http.StatusText(http.StatusMethodNotAllowed)
	if !bytes.Contains(w.Body.Bytes(), []byte(expected)) {
		t.Errorf("expected %q response, got %q", expected, w.Body.String())
	}
}
