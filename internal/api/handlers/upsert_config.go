package handlers

import (
	"encoding/json"
	"net/http"

	"example.com/myapp/internal/services"
)

// UpsertConfig handles HTTP requests to create or update a user's proxy chaos configuration.
// It expects a JSON body matching the services.Config structure and an "X-User-ID" header.
// It returns a JSON response containing the newly updated configuration.
func UpsertConfig(store services.ConfigStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userID := r.Header.Get("X-User-ID")
		if userID == "" {
			http.Error(w, "missing X-User-ID header", http.StatusBadRequest)
			return
		}

		// pull the configs from the request
		var config services.Config
		err := json.NewDecoder(r.Body).Decode(&config)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		new_config, ok := store.UpsertConfig(userID, config)
		if !ok {
			http.Error(w, "failed to update config", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(new_config)
	}
}
