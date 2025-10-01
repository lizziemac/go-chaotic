// Package api defines the application's HTTP layer.
//
// It contains route registration, request handlers, and response logic.
// This layer translates HTTP requests into calls to underlying services
// and returns appropriate HTTP responses.
//
// Typical responsibilities:
//   - Defining routes and route groups
//   - Binding and validating request payloads
//   - Writing JSON or other responses
//
// The api package should be minimal and delegate business logic
// to internal services.
package api
