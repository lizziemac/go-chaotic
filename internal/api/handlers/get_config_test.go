package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/myapp/internal/services"
)

func TestGetConfig(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		mockGetConfig  *services.Config
		expectedStatus int
	}{
		{
			name:           "successful get",
			userID:         "user-1",
			mockGetConfig:  &services.Config{Mode: 2},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "missing user id",
			userID:         "",
			mockGetConfig:  &services.Config{Mode: 2},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "not found",
			userID:         "user-2",
			mockGetConfig:  nil,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &mockStore{getConfig: tt.mockGetConfig}
			handler := GetConfig(store)

			req := httptest.NewRequest(http.MethodGet, "/proxy/api/v1/config", nil)
			if tt.userID != "" {
				req.Header.Set("X-User-ID", tt.userID)
			}
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
