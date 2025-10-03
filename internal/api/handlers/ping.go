package handlers

import (
	"encoding/json"
	"net/http"

	"example.com/myapp/internal/services"
)

type PingResponse struct {
	Message string `json:"message"`
}

// Ping handles HTTP GET requests to respond with a pong message.
// It returns a JSON response containing the message.
// If the request method is not GET, it responds with "method not allowed".
func Ping(w http.ResponseWriter, r *http.Request) {
	resp := PingResponse{Message: services.Pong()}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
