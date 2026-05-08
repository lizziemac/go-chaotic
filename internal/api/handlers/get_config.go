package handlers

import (
	"encoding/json"
	"net/http"

	"example.com/myapp/internal/services"
)

// GetConfig handles HTTP GET requests to retrieve a user's current chaos configuration.
// It expects an "X-User-ID" header.
func GetConfig(store services.ConfigStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-User-ID")
		if userID == "" {
			http.Error(w, "missing X-User-ID header", http.StatusBadRequest)
			return
		}

		config, ok := store.GetConfig(userID)
		if !ok {
			http.Error(w, "config not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(config)
	}
}
