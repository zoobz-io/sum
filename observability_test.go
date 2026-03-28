//go:build testing

package sum

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/zoobz-io/aperture"
)

func TestWithObservability(t *testing.T) {
	Reset()
	t.Cleanup(Reset)

	// Start a mock OTLP collector that accepts but discards data.
	collector := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer collector.Close()

	s := New()
	ctx := context.Background()

	err := s.WithObservability(ctx, "test-service", collector.URL)
	if err != nil {
		t.Fatalf("WithObservability failed: %v", err)
	}

	if s.Observability() == nil {
		t.Error("expected non-nil aperture instance")
	}

	s.mu.RLock()
	hasProviders := s.providers != nil
	s.mu.RUnlock()
	if !hasProviders {
		t.Error("expected non-nil OTEL providers")
	}
}

func TestObservabilityNilWithoutSetup(t *testing.T) {
	Reset()
	t.Cleanup(Reset)

	s := New()
	if s.Observability() != nil {
		t.Error("expected nil aperture before WithObservability")
	}
}

func TestWithObservabilityCustomSchema(t *testing.T) {
	Reset()
	t.Cleanup(Reset)

	collector := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer collector.Close()

	s := New()
	ctx := context.Background()

	err := s.WithObservability(ctx, "test-service", collector.URL)
	if err != nil {
		t.Fatalf("WithObservability failed: %v", err)
	}

	// Apply a custom schema with log whitelist.
	schema := aperture.Schema{
		Logs: &aperture.LogSchema{
			Whitelist: []string{"sum.connected"},
		},
	}
	err = s.Observability().Apply(schema)
	if err != nil {
		t.Fatalf("Apply schema failed: %v", err)
	}
}

func TestShutdownWithObservability(t *testing.T) {
	Reset()
	t.Cleanup(Reset)

	collector := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer collector.Close()

	s := New()
	ctx := context.Background()

	err := s.WithObservability(ctx, "test-service", collector.URL)
	if err != nil {
		t.Fatalf("WithObservability failed: %v", err)
	}

	// Start and shutdown to verify the three-phase ordering works.
	go func() {
		s.Start("127.0.0.1", 0)
	}()
	time.Sleep(100 * time.Millisecond)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(shutdownCtx); err != nil {
		t.Fatalf("Shutdown failed: %v", err)
	}
}

func TestShutdownWithoutObservability(t *testing.T) {
	Reset()
	t.Cleanup(Reset)

	s := New()

	go func() {
		s.Start("127.0.0.1", 0)
	}()
	time.Sleep(100 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		t.Fatalf("Shutdown failed: %v", err)
	}
}
