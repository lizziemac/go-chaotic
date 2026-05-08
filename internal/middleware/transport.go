package middleware

import (
	"errors"
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
	if userID == "" {
		return t.Next.RoundTrip(req)
	}

	config, ok := t.Config.GetConfig(userID)

	if ok {
		if config.Mode&services.DropMsg != 0 {
			var dropRate float32
			if config.DropRate != nil {
				dropRate = *config.DropRate
			}
			if rand.Float32() < dropRate {
				return nil, errors.New("connection dropped")
			}
		}
		if config.Mode&services.DelayMsg != 0 {
			var delay time.Duration
			if config.LatencyDelay != nil {
				delay = *config.LatencyDelay
			}
			time.Sleep(delay) // delayed!
		}
	}

	return t.Next.RoundTrip(req)
}
