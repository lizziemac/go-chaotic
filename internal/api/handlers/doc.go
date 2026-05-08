// Package handlers defines HTTP handler functions for the API routes.
//
// Handlers are responsible for:
//   - Decoding and validating requests
//   - Calling the appropriate service logic
//   - Writing responses and handling errors
//
// Each handler should be small and focused, delegating heavy logic
// to internal/services. Handlers map directly to API endpoints.
package handlers
