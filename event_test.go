//go:build testing

package sum

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/zoobz-io/capitan"
)

func init() {
	// Enable synchronous mode for deterministic event processing in tests
	capitan.Configure(capitan.WithSyncMode())
}

// testEventData is a struct type for testing Event with sentinel.
type testEventData struct {
	Message string
}

// testEventInt is a struct type for testing Event with int-like data.
type testEventInt struct {
	Value int
}

func TestNewEvent(t *testing.T) {
	t.Parallel()

	signal := capitan.NewSignal("test.event", "Test event")
	event := NewEvent[testEventData](signal, capitan.SeverityInfo)

	if event.Signal != signal {
		t.Errorf("expected signal %v, got %v", signal, event.Signal)
	}
	if event.level != capitan.SeverityInfo {
		t.Errorf("expected level %v, got %v", capitan.SeverityInfo, event.level)
	}
}

func TestNewEventConstructors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		factory  func(capitan.Signal) Event[testEventData]
		expected capitan.Severity
	}{
		{"debug", NewDebugEvent[testEventData], capitan.SeverityDebug},
		{"info", NewInfoEvent[testEventData], capitan.SeverityInfo},
		{"warn", NewWarnEvent[testEventData], capitan.SeverityWarn},
		{"error", NewErrorEvent[testEventData], capitan.SeverityError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			signal := capitan.NewSignal("test."+tt.name, "Test "+tt.name+" event")
			event := tt.factory(signal)

			if event.level != tt.expected {
				t.Errorf("expected level %v, got %v", tt.expected, event.level)
			}
		})
	}
}

func TestEventEmitAndListen(t *testing.T) {
	t.Parallel()

	signal := capitan.NewSignal("test.emit.listen", "Emit and listen test")
	event := NewInfoEvent[testEventData](signal)

	var received testEventData
	var mu sync.Mutex
	done := make(chan struct{})

	listener := event.Listen(func(ctx context.Context, data testEventData) {
		mu.Lock()
		received = data
		mu.Unlock()
		close(done)
	})
	defer listener.Close()

	ctx := context.Background()
	event.Emit(ctx, testEventData{Message: "hello"})

	select {
	case <-done:
		mu.Lock()
		if received.Message != "hello" {
			t.Errorf("expected 'hello', got '%s'", received.Message)
		}
		mu.Unlock()
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for event")
	}
}

func TestEventListenOnce(t *testing.T) {
	t.Parallel()

	signal := capitan.NewSignal("test.listen.once", "Listen once test")
	event := NewInfoEvent[testEventInt](signal)

	var count int
	var mu sync.Mutex
	done := make(chan struct{}, 2)

	listener := event.ListenOnce(func(ctx context.Context, data testEventInt) {
		mu.Lock()
		count++
		mu.Unlock()
		done <- struct{}{}
	})
	defer listener.Close()

	ctx := context.Background()
	event.Emit(ctx, testEventInt{Value: 1})
	event.Emit(ctx, testEventInt{Value: 2})

	// Wait for potential callbacks
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for first event")
	}

	// Give time for second event to process (it shouldn't trigger callback)
	time.Sleep(50 * time.Millisecond)

	mu.Lock()
	if count != 1 {
		t.Errorf("expected callback count 1, got %d", count)
	}
	mu.Unlock()
}

func TestEventEmitAllSeverities(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		level   capitan.Severity
		factory func(capitan.Signal) Event[testEventData]
	}{
		{"debug", capitan.SeverityDebug, NewDebugEvent[testEventData]},
		{"info", capitan.SeverityInfo, NewInfoEvent[testEventData]},
		{"warn", capitan.SeverityWarn, NewWarnEvent[testEventData]},
		{"error", capitan.SeverityError, NewErrorEvent[testEventData]},
		{"default", capitan.Severity("CUSTOM"), func(s capitan.Signal) Event[testEventData] {
			return NewEvent[testEventData](s, capitan.Severity("CUSTOM"))
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			signal := capitan.NewSignal("test.severity."+tt.name, "Severity "+tt.name+" test")
			event := tt.factory(signal)

			var received bool
			done := make(chan struct{})

			listener := event.Listen(func(ctx context.Context, data testEventData) {
				received = true
				close(done)
			})
			defer listener.Close()

			ctx := context.Background()
			event.Emit(ctx, testEventData{Message: "test"})

			select {
			case <-done:
				if !received {
					t.Error("expected to receive event")
				}
			case <-time.After(time.Second):
				t.Fatal("timeout waiting for event")
			}
		})
	}
}
