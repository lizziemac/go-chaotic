// Package services contains the core business logic of the application.
//
// Services coordinate lower-level operations (database, APIs, etc.)
// and implement use-case-specific behavior. They should not depend
// on HTTP-specific details.
//
// Typical responsibilities:
//   - Application workflows
//   - Business rules and validation
//   - Interfacing with repositories or external systems
//
// Keep services testable and independent from framework code.
package services
