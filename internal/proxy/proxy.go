package reverseproxy

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"example.com/myapp/internal/logger"
	"example.com/myapp/internal/middleware"
	"example.com/myapp/internal/services"
)

// Set up the reverse proxy with the target URL that the proxy will be in front of. This can be set with the
// TARGET_URL environment variable
func Setup(ctx context.Context, configStore services.ConfigStore) *httputil.ReverseProxy {
	// Get the URL that the reverse proxy will be in front of
	var targetURL = os.Getenv("TARGET_URL")
	target, err := url.Parse(targetURL)
	if err != nil || targetURL == "" {
		logger.Fatal("Invalid target URL", "error", err)
		os.Exit(1)
	}

	// Set up the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Transport = &middleware.Transport{
		Next:   http.DefaultTransport,
		Logger: logger.Logger,
		Config: configStore,
	}

	return proxy
}
