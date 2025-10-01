// Package middleware provides HTTP middleware used within the API layer.
//
// Middleware functions wrap HTTP handlers to aid in the API's cross-cutting concerns, such as:
//   - Authz/Authn
//   - Rate limiting / quota enforcement
//   - CORS handling
//   - Business specific logging
//
// Each middleware function is typically defined as:
//
//	func RequestTimer(next StateHandler) StateHandler
//
// where StateHandler is a custom handler type that passes shared state between middleware layers.
package middleware
