package middleware

import (
	"log/slog"
	"math/rand/v2"
	"net/http"
	"time"

	"example.com/myapp/internal/services"
)

type Transport struct {
	Next   http.RoundTripper
	Logger *slog.Logger
	Config services.ConfigStore
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	userID := req.Header.Get("X-User-ID")
	t.Logger.Info("INJECTION")
	if userID == "" {
		return t.Next.RoundTrip(req)
	}

	config, ok := t.Config.GetConfig(userID)

	if !ok || config.Mode == services.PassMsg {
		return t.Next.RoundTrip(req)
	}

	var delay time.Duration
	if config.Mode&services.DropMsg != 0 {
		if rand.Float32() < config.DropRate {
			return nil, nil // dropped!
		}
	}
	if config.Mode&services.DelayMsg != 0 {
		time.Sleep(delay) // delayed!
	}

	return t.Next.RoundTrip(req)
}
