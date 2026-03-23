package sum

import (
	"context"

	"github.com/zoobz-io/capitan"
	"github.com/zoobz-io/sentinel"
)

type (
	// Signal is a named event identifier.
	Signal = capitan.Signal
	// Severity is the importance level of an emission.
	Severity = capitan.Severity
)

// NewSignal creates a signal with name and description.
var NewSignal = capitan.NewSignal

// Emission function re-exports.
var (
	Emit  = capitan.Emit
	Debug = capitan.Debug
	Info  = capitan.Info
	Warn  = capitan.Warn
	Error = capitan.Error
)

// Field key constructor re-exports.
var (
	NewStringKey   = capitan.NewStringKey
	NewIntKey      = capitan.NewIntKey
	NewInt64Key    = capitan.NewInt64Key
	NewFloat64Key  = capitan.NewFloat64Key
	NewBoolKey     = capitan.NewBoolKey
	NewTimeKey     = capitan.NewTimeKey
	NewDurationKey = capitan.NewDurationKey
	NewErrorKey    = capitan.NewErrorKey
)

// Event provides bidirectional access to a signal with typed data.
// Use Emit to dispatch events and Listen to register callbacks.
type Event[T any] struct {
	// Signal is the underlying capitan signal.
	Signal capitan.Signal

	// Key is the typed key for extracting event data.
	Key capitan.GenericKey[T]

	// level determines the severity used when emitting.
	level capitan.Severity
}

// Emit dispatches an event with the configured severity level.
func (e Event[T]) Emit(ctx context.Context, data T) {
	switch e.level {
	case capitan.SeverityDebug:
		capitan.Debug(ctx, e.Signal, e.Key.Field(data))
	case capitan.SeverityInfo:
		capitan.Info(ctx, e.Signal, e.Key.Field(data))
	case capitan.SeverityWarn:
		capitan.Warn(ctx, e.Signal, e.Key.Field(data))
	case capitan.SeverityError:
		capitan.Error(ctx, e.Signal, e.Key.Field(data))
	default:
		capitan.Emit(ctx, e.Signal, e.Key.Field(data))
	}
}

// Listen registers a callback for this event.
// Returns a Listener that can be closed to unregister.
func (e Event[T]) Listen(callback func(context.Context, T)) *capitan.Listener {
	return capitan.Hook(e.Signal, func(ctx context.Context, ev *capitan.Event) {
		if data, ok := e.Key.From(ev); ok {
			callback(ctx, data)
		}
	})
}

// ListenOnce registers a callback that fires only once, then automatically unregisters.
// Returns a Listener that can be closed early to prevent the callback from firing.
func (e Event[T]) ListenOnce(callback func(context.Context, T)) *capitan.Listener {
	return capitan.HookOnce(e.Signal, func(ctx context.Context, ev *capitan.Event) {
		if data, ok := e.Key.From(ev); ok {
			callback(ctx, data)
		}
	})
}

// NewEvent creates an Event with the given signal and severity level.
// The variant is derived automatically from T via sentinel.
func NewEvent[T any](signal capitan.Signal, level capitan.Severity) Event[T] {
	meta := sentinel.Inspect[T]()
	return Event[T]{
		Signal: signal,
		Key:    capitan.NewKey[T]("data", capitan.Variant(meta.FQDN)),
		level:  level,
	}
}

// NewDebugEvent creates an Event that emits at Debug level.
func NewDebugEvent[T any](signal capitan.Signal) Event[T] {
	return NewEvent[T](signal, capitan.SeverityDebug)
}

// NewInfoEvent creates an Event that emits at Info level.
func NewInfoEvent[T any](signal capitan.Signal) Event[T] {
	return NewEvent[T](signal, capitan.SeverityInfo)
}

// NewWarnEvent creates an Event that emits at Warn level.
func NewWarnEvent[T any](signal capitan.Signal) Event[T] {
	return NewEvent[T](signal, capitan.SeverityWarn)
}

// NewErrorEvent creates an Event that emits at Error level.
func NewErrorEvent[T any](signal capitan.Signal) Event[T] {
	return NewEvent[T](signal, capitan.SeverityError)
}
