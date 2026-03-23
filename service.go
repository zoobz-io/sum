package sum

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/zoobz-io/cereal"
	"github.com/zoobz-io/rocco"
	"github.com/zoobz-io/scio"
)

// service is the singleton instance.
var (
	instance *Service
	once     sync.Once
)

// Service wraps a rocco engine and scio catalog, providing application lifecycle.
type Service struct {
	encryptors map[cereal.EncryptAlgo]cereal.Encryptor
	hashers    map[cereal.HashAlgo]cereal.Hasher
	maskers    map[cereal.MaskType]cereal.Masker
	engine     *rocco.Engine
	catalog    *scio.Scio
	codec      cereal.Codec
	mu         sync.RWMutex
}

// New creates or returns the singleton Service.
// Subsequent calls return the existing instance.
func New() *Service {
	once.Do(func() {
		instance = &Service{
			engine:     rocco.NewEngine(),
			catalog:    scio.New(),
			encryptors: make(map[cereal.EncryptAlgo]cereal.Encryptor),
			hashers:    make(map[cereal.HashAlgo]cereal.Hasher),
			maskers:    make(map[cereal.MaskType]cereal.Masker),
		}
	})
	return instance
}

// svc returns the singleton, panicking if not initialized.
func svc() *Service {
	if instance == nil {
		panic("sum: service not initialized, call New() first")
	}
	return instance
}

// Handle registers endpoints with the underlying engine.
func (s *Service) Handle(endpoints ...rocco.Endpoint) {
	s.engine.WithHandlers(endpoints...)
}

// Tag registers an OpenAPI tag with a description.
func (s *Service) Tag(name, description string) {
	s.engine.WithTag(name, description)
}

// Engine returns the underlying rocco engine for advanced usage.
func (s *Service) Engine() *rocco.Engine {
	return s.engine
}

// Catalog returns the scio data catalog for advanced usage.
func (s *Service) Catalog() *scio.Scio {
	return s.catalog
}

// WithEncryptor registers an encryptor for the given algorithm.
func (s *Service) WithEncryptor(algo cereal.EncryptAlgo, enc cereal.Encryptor) *Service {
	s.mu.Lock()
	s.encryptors[algo] = enc
	s.mu.Unlock()
	return s
}

// WithHasher registers a hasher for the given algorithm.
func (s *Service) WithHasher(algo cereal.HashAlgo, h cereal.Hasher) *Service {
	s.mu.Lock()
	s.hashers[algo] = h
	s.mu.Unlock()
	return s
}

// WithMasker registers a masker for the given mask type.
func (s *Service) WithMasker(mt cereal.MaskType, m cereal.Masker) *Service {
	s.mu.Lock()
	s.maskers[mt] = m
	s.mu.Unlock()
	return s
}

// WithCodec sets the default codec for cereal processors and the rocco engine.
func (s *Service) WithCodec(codec cereal.Codec) *Service {
	s.mu.Lock()
	s.codec = codec
	s.mu.Unlock()
	s.engine.WithCodec(roccoCodec{codec})
	return s
}

// Start begins serving. This method blocks until shutdown.
func (s *Service) Start(host string, port int) error {
	return s.engine.Start(host, port)
}

// Shutdown gracefully stops the service.
func (s *Service) Shutdown(ctx context.Context) error {
	if s.engine == nil {
		return fmt.Errorf("service not started")
	}
	return s.engine.Shutdown(ctx)
}

// Run starts the service and blocks until a shutdown signal is received.
// Handles SIGINT and SIGTERM, then performs graceful shutdown with a 30 second timeout.
func (s *Service) Run(host string, port int) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- s.engine.Start(host, port)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()
		return s.Shutdown(shutdownCtx)
	}
}
