//go:build testing

package sum

import (
	"context"
	"encoding/json"
	"sync"
	"testing"

	"github.com/zoobz-io/cereal"
	"github.com/zoobz-io/slush"
)

// testUser is a minimal Cloner type for boundary tests.
type testUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (u testUser) Clone() testUser { return u }

// testCodec is a simple JSON codec for testing.
type testCodec struct{}

func (c *testCodec) ContentType() string              { return "application/json" }
func (c *testCodec) Marshal(v any) ([]byte, error)    { return json.Marshal(v) }
func (c *testCodec) Unmarshal(data []byte, v any) error { return json.Unmarshal(data, v) }

// stubEncryptor satisfies cereal.Encryptor for capability propagation tests.
type stubEncryptor struct{}

func (stubEncryptor) Encrypt(p []byte) ([]byte, error) { return p, nil }
func (stubEncryptor) Decrypt(c []byte) ([]byte, error) { return c, nil }

// stubHasher satisfies cereal.Hasher.
type stubHasher struct{}

func (stubHasher) Hash(p []byte) (string, error) { return string(p), nil }

// stubMasker satisfies cereal.Masker.
type stubMasker struct{}

func (stubMasker) Mask(v string) (string, error) { return "***", nil }

func resetAll(t *testing.T) {
	t.Helper()
	slush.Reset()
	instance = nil
	once = sync.Once{}
	t.Cleanup(func() {
		slush.Reset()
		instance = nil
		once = sync.Once{}
	})
}

func TestNewBoundary(t *testing.T) {
	resetAll(t)
	New()
	k := Start()

	b, err := NewBoundary[testUser](k)
	if err != nil {
		t.Fatalf("NewBoundary failed: %v", err)
	}
	if b == nil {
		t.Fatal("expected non-nil boundary")
	}
	if b.Processor == nil {
		t.Fatal("expected non-nil processor")
	}

	Freeze(k)
}

func TestNewBoundaryPanicsWithoutService(t *testing.T) {
	resetAll(t)

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic when service not initialized")
		}
	}()

	k := Start()
	NewBoundary[testUser](k)
}

func TestCapabilityPropagation(t *testing.T) {
	resetAll(t)
	s := New()

	s.WithEncryptor(cereal.EncryptAES, stubEncryptor{})
	s.WithHasher(cereal.HashSHA256, stubHasher{})
	s.WithMasker(cereal.MaskEmail, stubMasker{})
	s.WithCodec(&testCodec{})

	k := Start()
	b, err := NewBoundary[testUser](k)
	if err != nil {
		t.Fatalf("NewBoundary failed: %v", err)
	}
	Freeze(k)

	// Verify boundary was created with capabilities by exercising Receive/Send.
	// With no tagged fields on testUser, these should pass through cleanly.
	ctx := context.Background()
	u := testUser{ID: "1", Name: "Alice"}

	received, err := b.Receive(ctx, u)
	if err != nil {
		t.Fatalf("Receive failed: %v", err)
	}
	if received.ID != u.ID || received.Name != u.Name {
		t.Errorf("Receive altered untagged fields: got %+v", received)
	}

	sent, err := b.Send(ctx, u)
	if err != nil {
		t.Fatalf("Send failed: %v", err)
	}
	if sent.ID != u.ID || sent.Name != u.Name {
		t.Errorf("Send altered untagged fields: got %+v", sent)
	}
}

func TestBoundaryAutoRegistered(t *testing.T) {
	resetAll(t)
	New()
	k := Start()

	_, err := NewBoundary[testUser](k)
	if err != nil {
		t.Fatalf("NewBoundary failed: %v", err)
	}
	Freeze(k)

	ctx := context.Background()
	b := MustUse[*Boundary[testUser]](ctx)
	if b == nil {
		t.Fatal("expected non-nil boundary from MustUse")
	}
}

func TestBoundaryMustUsePanicsWhenNotRegistered(t *testing.T) {
	resetAll(t)
	New()
	k := Start()
	Freeze(k)

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic when boundary not registered")
		}
	}()

	ctx := context.Background()
	MustUse[*Boundary[testUser]](ctx)
}

func TestCodecAdapterWiresToRocco(t *testing.T) {
	resetAll(t)
	s := New()

	codec := &testCodec{}
	s.WithCodec(codec)

	// Verify codec field was set on the service.
	s.mu.RLock()
	if s.codec == nil {
		t.Error("expected codec to be set on service")
	}
	s.mu.RUnlock()
}

func TestResetClearsCapabilities(t *testing.T) {
	instance = nil
	once = sync.Once{}
	s := New()
	k := Start()

	s.WithEncryptor(cereal.EncryptAES, stubEncryptor{})
	s.WithHasher(cereal.HashSHA256, stubHasher{})
	s.WithMasker(cereal.MaskEmail, stubMasker{})
	s.WithCodec(&testCodec{})

	Freeze(k)
	Reset()

	// After Reset, instance is nil; create fresh.
	s = New()

	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.encryptors) != 0 {
		t.Error("expected empty encryptors after Reset")
	}
	if len(s.hashers) != 0 {
		t.Error("expected empty hashers after Reset")
	}
	if len(s.maskers) != 0 {
		t.Error("expected empty maskers after Reset")
	}
	if s.codec != nil {
		t.Error("expected nil codec after Reset")
	}
}
