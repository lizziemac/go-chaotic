// Package logger provides structured logging utilities for the application.
//
// It centralizes logging configuration and offers helpers for consistent
// log formatting, context propagation, and log levels.
//
// Typical responsibilities:
//   - Initialize global or context-aware loggers
//   - Provide standardized log methods
//   - Integrate with observability stacks (Zap, Zerolog, etc.)
//
// The logger package should be imported by other internal packages
// that require structured logging.
package logger
