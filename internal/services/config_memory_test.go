package services

import (
	"sync"
	"testing"
	"time"
)

func TestUpsertConfig_Success(t *testing.T) {
	ctx := t.Context()
	registry := NewConfigRegistry(ctx)

	cfg := Config{Mode: DelayMsg}
	inserted, ok := registry.UpsertConfig("user-1", cfg)

	if !ok {
		t.Fatal("expected upsert to succeed")
	}
	if inserted.Mode != DelayMsg {
		t.Errorf("expected mode %v, got %v", DelayMsg, inserted.Mode)
	}
	if inserted.TTL.IsZero() {
		t.Error("expected TTL to be set on upsert")
	}
}

func TestGetConfig_Exists(t *testing.T) {
	ctx := t.Context()
	registry := NewConfigRegistry(ctx)

	// Seed the registry
	registry.UpsertConfig("user-1", Config{Mode: DropMsg})

	cfg, ok := registry.GetConfig("user-1")
	if !ok {
		t.Fatal("expected to find config")
	}
	if cfg.Mode != DropMsg {
		t.Errorf("expected mode %v, got %v", DropMsg, cfg.Mode)
	}
}

func TestGetConfig_NotFound(t *testing.T) {
	ctx := t.Context()
	registry := NewConfigRegistry(ctx)

	_, ok := registry.GetConfig("unknown-user")
	if ok {
		t.Fatal("expected config to not be found")
	}
}

func TestCleanup_RemovesExpiredConfigs(t *testing.T) {
	ctx := t.Context()
	registry := NewConfigRegistry(ctx)

	// Insert configs manually to bypass the TTL enforcement in UpsertConfig
	registry.mu.Lock()
	registry.configs["expired-user"] = Config{
		Mode: DropMsg,
		TTL:  time.Now().Add(-1 * time.Hour), // Expired 1 hour ago
	}
	registry.configs["valid-user"] = Config{
		Mode: PassMsg,
		TTL:  time.Now().Add(1 * time.Hour), // Expires in 1 hour
	}
	registry.mu.Unlock()

	// Manually trigger the cleanup routine
	registry.cleanup()

	if _, ok := registry.GetConfig("expired-user"); ok {
		t.Error("expected expired-user to be cleaned up")
	}

	if _, ok := registry.GetConfig("valid-user"); !ok {
		t.Error("expected valid-user to remain in registry")
	}
}

func TestConcurrentAccess(t *testing.T) {
	ctx := t.Context()
	registry := NewConfigRegistry(ctx)

	var wg sync.WaitGroup
	userID := "concurrent-user"

	// Spawn multiple goroutines to Upsert, Get, and Cleanup concurrently
	// to ensure there are no race conditions or deadlocks.
	for range 100 {
		wg.Go(func() {
			registry.UpsertConfig(userID, Config{Mode: DelayMsg})
		})
		wg.Go(func() {
			registry.GetConfig(userID)
		})
		wg.Go(func() {
			registry.cleanup()
		})
	}

	wg.Wait()
}
