// Package middleware provides reusable HTTP middleware functions.
//
// Middleware wraps handlers to perform cross-cutting tasks such as:
//   - Logging
//   - Authentication and authorization
//   - CORS handling
//   - Request tracing
//   - Panic recovery
//
// Each middleware is typically implemented as:
//
//	func RequestTimer(next StateHandler) StateHandler
//
// where StateHandler is a custom handler type that passes shared state between middleware.
package middleware
