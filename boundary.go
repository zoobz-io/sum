package sum

import (
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

// Boundary wraps a cereal Processor and auto-registers with the service registry.
type Boundary[T cereal.Cloner[T]] struct {
	*cereal.Processor[T]
}

// NewBoundary creates a Boundary[T], applies shared capabilities from the Service,
// and registers it in the service registry under the given key.
func NewBoundary[T cereal.Cloner[T]](k Key) (*Boundary[T], error) {
	s := svc()
	proc, err := cereal.NewProcessor[T]()
	if err != nil {
		return nil, err
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

	b := &Boundary[T]{Processor: proc}
	Register[*Boundary[T]](k, b)
	return b, nil
}

// roccoCodec adapts a cereal.Codec to rocco.Codec.
type roccoCodec struct{ cereal.Codec }

var _ rocco.Codec = roccoCodec{}

func (r roccoCodec) ContentType() string                { return r.Codec.ContentType() }
func (r roccoCodec) Marshal(v any) ([]byte, error)      { return r.Codec.Marshal(v) }
func (r roccoCodec) Unmarshal(data []byte, v any) error { return r.Codec.Unmarshal(data, v) }
