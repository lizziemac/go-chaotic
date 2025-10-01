// Package internal contains non-exported, application-specific code.
//
// Code inside internal is private to this module and cannot be imported
// by other modules. It is used to group implementation details that
// should not form part of the public API.
//
// Typical contents include:
//   - Core business logic
//   - Middleware
//   - Internal APIs or handlers
//   - Utility packages (e.g., config, logger)
//
// Example structure:
//
//	internal/
//	  api/
//	  middleware/
//	  services/
//	  logger/
//
// External code should import only public packages outside internal
package internal
