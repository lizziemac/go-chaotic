package services

import (
	"time"
)

// A bitmapped value to indicate how to handle the message
type MessageMode int

const (
	PassMsg  MessageMode = 0      // do nothing to the message
	DropMsg  MessageMode = 1 << 0 // lose a certain amount of incoming messages (binary 01)
	DelayMsg MessageMode = 1 << 1 // hold the message for a certain amount of time (binary 10)
)

const TICKER_INTERVAL = 4 * time.Hour

const TTL_DAYS = 30
const MAX_TTL = TTL_DAYS * (24 * time.Hour)

// Holds the current proxy rules
type Config struct {
	Mode         MessageMode    `json:"mode"`             // What chaos will be applied to this message
	DropRate     *float32       `json:"drop_rate"`        // If message mode is 'drop', the percentage rate at which messages are dropped
	LatencyDelay *time.Duration `json:"latency_delay_ns"` // If message mode is 'delay', the amount of time it takes for each message to be forwarded
	TTL          time.Time      `json:"ttl"`              // When to clear the user key, if it hasn't been updated in a while
}

// ConfigStore defines the contract for fetching and updating configurations.
type ConfigStore interface {
	UpsertConfig(userID string, config Config) (*Config, bool)
	GetConfig(userID string) (*Config, bool)
}
