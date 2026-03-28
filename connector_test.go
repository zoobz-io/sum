//go:build testing

package sum

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/zoobz-io/capitan"
)

type mockClient struct {
	closed atomic.Bool
}

func (m *mockClient) Close() error {
	m.closed.Store(true)
	return nil
}

type mockClientNoCloser struct {
	value string
}

func TestConnect(t *testing.T) {
	Reset()
	t.Cleanup(Reset)
	New()
	k := Start()
	ctx := context.Background()

	client := &mockClient{}
	err := Connect[*mockClient](ctx, k, "test-db", func(ctx context.Context) (*mockClient, error) {
		return client, nil
	})
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	got, err := Use[*mockClient](ctx)
	if err != nil {
		t.Fatalf("Use failed: %v", err)
	}
	if got != client {
		t.Error("expected same client instance from registry")
	}
}

func TestConnectFactoryError(t *testing.T) {
	Reset()
	t.Cleanup(Reset)
	New()
	k := Start()
	ctx := context.Background()

	err := Connect[*mockClient](ctx, k, "bad-db", func(ctx context.Context) (*mockClient, error) {
		return nil, errors.New("connection refused")
	})
	if err == nil {
		t.Fatal("expected error from factory")
	}
	if got := err.Error(); got != "connect bad-db: connection refused" {
		t.Errorf("unexpected error: %s", got)
	}
}

func TestConnectTracksCloser(t *testing.T) {
	Reset()
	t.Cleanup(Reset)
	s := New()
	k := Start()
	ctx := context.Background()

	err := Connect[*mockClient](ctx, k, "closeable", func(ctx context.Context) (*mockClient, error) {
		return &mockClient{}, nil
	})
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	s.mu.RLock()
	n := len(s.closers)
	s.mu.RUnlock()
	if n != 1 {
		t.Errorf("expected 1 closer, got %d", n)
	}
}

func TestConnectNoCloserNotTracked(t *testing.T) {
	Reset()
	t.Cleanup(Reset)
	s := New()
	k := Start()
	ctx := context.Background()

	err := Connect[*mockClientNoCloser](ctx, k, "no-closer", func(ctx context.Context) (*mockClientNoCloser, error) {
		return &mockClientNoCloser{value: "test"}, nil
	})
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	s.mu.RLock()
	n := len(s.closers)
	s.mu.RUnlock()
	if n != 0 {
		t.Errorf("expected 0 closers, got %d", n)
	}
}

func TestConnectEmitsSignal(t *testing.T) {
	Reset()
	t.Cleanup(Reset)
	New()
	k := Start()
	ctx := context.Background()

	var received bool
	listener := capitan.Hook(SignalConnected, func(ctx context.Context, ev *capitan.Event) {
		name, ok := KeyConnectorName.From(ev)
		if ok && name == "signal-test" {
			received = true
		}
	})
	defer listener.Close()

	err := Connect[*mockClient](ctx, k, "signal-test", func(ctx context.Context) (*mockClient, error) {
		return &mockClient{}, nil
	})
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	if !received {
		t.Error("expected SignalConnected to be emitted")
	}
}

func TestShutdownClosesInReverseOrder(t *testing.T) {
	Reset()
	t.Cleanup(Reset)
	s := New()
	k := Start()
	ctx := context.Background()

	first := &mockClient{}
	Connect[*mockClient](ctx, k, "first", func(ctx context.Context) (*mockClient, error) {
		return first, nil
	})

	second := &mockClient{}
	third := &mockClient{}
	s.mu.Lock()
	s.closers = append(s.closers,
		namedCloser{name: "second", closer: second},
		namedCloser{name: "third", closer: third},
	)
	s.mu.Unlock()

	var closeOrder []string
	listener := capitan.Hook(SignalDisconnected, func(ctx context.Context, ev *capitan.Event) {
		if name, ok := KeyConnectorName.From(ev); ok {
			closeOrder = append(closeOrder, name)
		}
	})
	defer listener.Close()

	// Start and shutdown after server is ready
	errCh := make(chan error, 1)
	go func() {
		errCh <- s.Start("127.0.0.1", 0)
	}()
	time.Sleep(100 * time.Millisecond)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(shutdownCtx); err != nil {
		t.Fatalf("Shutdown failed: %v", err)
	}

	if len(closeOrder) != 3 {
		t.Fatalf("expected 3 disconnects, got %d: %v", len(closeOrder), closeOrder)
	}
	if closeOrder[0] != "third" || closeOrder[1] != "second" || closeOrder[2] != "first" {
		t.Errorf("expected reverse order [third second first], got %v", closeOrder)
	}

	if !first.closed.Load() || !second.closed.Load() || !third.closed.Load() {
		t.Error("expected all clients to be closed")
	}
}

func TestShutdownNoClosers(t *testing.T) {
	Reset()
	t.Cleanup(Reset)
	s := New()

	errCh := make(chan error, 1)
	go func() {
		errCh <- s.Start("127.0.0.1", 0)
	}()
	time.Sleep(100 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		t.Fatalf("Shutdown failed: %v", err)
	}

	select {
	case err := <-errCh:
		if err != nil && err.Error() != "http: Server closed" {
			t.Errorf("unexpected start error: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Error("timeout waiting for server to stop")
	}
}
