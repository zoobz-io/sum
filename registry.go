package sum

import (
	"context"

	"github.com/zoobz-io/capitan"
	"github.com/zoobz-io/slush"
)

type (
	// Guard is a validation function that permits or denies service access.
	Guard = slush.Guard
	// ServiceInfo describes a registered service for enumeration.
	ServiceInfo = slush.ServiceInfo
	// Key grants the capability to register services.
	Key = slush.Key
)

// Handle configures a registered service with optional guards.
// Wraps slush.Handle to provide sum-specific conveniences.
type Handle[T any] struct {
	*slush.Handle[T]
}

// Guard adds a custom guard function to the service.
// Returns the Handle for chaining.
func (h *Handle[T]) Guard(g Guard) *Handle[T] {
	h.Handle.Guard(g)
	return h
}

// For restricts service access to contexts bearing any of the provided tokens.
// Equivalent to Guard(Require(tokens...)).
func (h *Handle[T]) For(tokens ...Token) *Handle[T] {
	h.Handle.Guard(Require(tokens...))
	return h
}

// Error re-exports from slush.
var (
	ErrNotFound     = slush.ErrNotFound
	ErrAccessDenied = slush.ErrAccessDenied
	ErrInvalidKey   = slush.ErrInvalidKey
)

// Signal re-exports from slush.
var (
	SignalRegistered capitan.Signal = slush.SignalRegistered
	SignalAccessed   capitan.Signal = slush.SignalAccessed
	SignalDenied     capitan.Signal = slush.SignalDenied
	SignalNotFound   capitan.Signal = slush.SignalNotFound
)

// Field key re-exports from slush.
var (
	KeyInterface capitan.Key = slush.KeyInterface
	KeyImpl      capitan.Key = slush.KeyImpl
	KeyError     capitan.Key = slush.KeyError
)

// Start initializes the service registry and returns a Key for registration.
// Panics if called more than once.
func Start() Key {
	return slush.Start()
}

// Freeze prevents further service registration.
// Panics if key is invalid.
func Freeze(k Key) {
	slush.Freeze(k)
}

// Register registers a service implementation for the contract type T.
// Returns a Handle for optional guard configuration.
// Panics if Start has not been called, key is invalid, or registry is frozen.
func Register[T any](k Key, impl T) *Handle[T] {
	return &Handle[T]{slush.Register[T](k, impl)}
}

// Use retrieves a service by its contract type T.
// Runs all registered guards with the provided context.
// Returns ErrNotFound if not registered, ErrAccessDenied if a guard fails.
func Use[T any](ctx context.Context) (T, error) {
	return slush.Use[T](ctx)
}

// MustUse retrieves a service by its contract type T.
// Panics if the service is not registered or a guard fails.
func MustUse[T any](ctx context.Context) T {
	return slush.MustUse[T](ctx)
}

// Services returns information about all registered services.
// Returns ErrInvalidKey if the key is invalid.
func Services(k Key) ([]ServiceInfo, error) {
	return slush.Services(k)
}
