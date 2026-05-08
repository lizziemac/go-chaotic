package services

import (
	"context"
	"sync"
	"time"
)

// Compile-time check to guarantee ConfigRegistry implements ConfigStore.
// If a method is missing or misspelled, the compiler will fail right here.
var _ ConfigStore = (*ConfigRegistry)(nil)

// The configuration registry object
type ConfigRegistry struct {
	mu      sync.RWMutex      // The mutext to protect the configs
	configs map[string]Config // The map of configs, by user
}

// NewConfigRegistry creates a new in-memory configuration store and starts the cleanup ticker.
func NewConfigRegistry(ctx context.Context) *ConfigRegistry {
	registry := &ConfigRegistry{
		configs: make(map[string]Config),
	}

	ticker := time.NewTicker(TICKER_INTERVAL)
	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				registry.cleanup()
			}
		}
	}()

	return registry
}

// Delete users from the config who have expired
func (r *ConfigRegistry) cleanup() {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	for userID, cfg := range r.configs {
		if now.After(cfg.TTL) {
			delete(r.configs, userID)
		}
	}
}

// Update the configuration for a specific user
func (r *ConfigRegistry) UpsertConfig(userID string, config Config) (*Config, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.configs == nil {
		r.configs = make(map[string]Config)
	}
	config.TTL = time.Now().Add(MAX_TTL)
	r.configs[userID] = config
	return &config, true
}

// Get the current configuration for a user
func (r *ConfigRegistry) GetConfig(userID string) (*Config, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.configs == nil {
		return nil, false
	}
	config, ok := r.configs[userID]
	if !ok {
		return nil, false
	}
	return &config, true
}
