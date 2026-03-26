package sum

import (
	"context"

	"github.com/zoobz-io/cereal"
	"github.com/zoobz-io/rocco"
)

type (
	// Codec is a re-export of cereal.Codec.
	Codec = cereal.Codec
	// Encryptor is a re-export of cereal.Encryptor.
	Encryptor = cereal.Encryptor
	// Hasher is a re-export of cereal.Hasher.
	Hasher = cereal.Hasher
	// Masker is a re-export of cereal.Masker.
	Masker = cereal.Masker
	// EncryptAlgo is a re-export of cereal.EncryptAlgo.
	EncryptAlgo = cereal.EncryptAlgo
	// HashAlgo is a re-export of cereal.HashAlgo.
	HashAlgo = cereal.HashAlgo
	// MaskType is a re-export of cereal.MaskType.
	MaskType = cereal.MaskType
)

// Boundary defines the serialization lifecycle operations for a type.
type Boundary[T any] interface {
	Send(ctx context.Context, obj T) (T, error)
	Receive(ctx context.Context, obj T) (T, error)
	Store(ctx context.Context, obj T) (T, error)
	Load(ctx context.Context, obj T) (T, error)
}

// boundary wraps a cereal Processor and auto-registers with the service registry.
type boundary[T cereal.Cloner[T]] struct {
	*cereal.Processor[T]
}

// NewBoundary creates a Boundary[T], applies shared capabilities from the Service,
// and registers it in the service registry under the given key.
// Panics if cereal.NewProcessor fails (structurally unreachable for valid Cloner types).
func NewBoundary[T cereal.Cloner[T]](k Key) Boundary[T] {
	s := svc()
	proc, err := cereal.NewProcessor[T]()
	if err != nil {
		panic("sum: NewBoundary: " + err.Error())
	}

	s.mu.RLock()
	for algo, enc := range s.encryptors {
		proc.SetEncryptor(algo, enc)
	}
	for algo, h := range s.hashers {
		proc.SetHasher(algo, h)
	}
	for mt, m := range s.maskers {
		proc.SetMasker(mt, m)
	}
	if s.codec != nil {
		proc.SetCodec(s.codec)
	}
	s.mu.RUnlock()

	b := &boundary[T]{Processor: proc}
	Register[Boundary[T]](k, b)
	return b
}

// roccoCodec adapts a cereal.Codec to rocco.Codec.
type roccoCodec struct{ cereal.Codec }

var _ rocco.Codec = roccoCodec{}

func (r roccoCodec) ContentType() string                { return r.Codec.ContentType() }
func (r roccoCodec) Marshal(v any) ([]byte, error)      { return r.Codec.Marshal(v) }
func (r roccoCodec) Unmarshal(data []byte, v any) error { return r.Codec.Unmarshal(data, v) }
