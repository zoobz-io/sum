//go:build testing || integration

package sum

import (
	"context"
	"sync"

	"github.com/zoobz-io/slush"
)

// Reset clears all registered services and resets initialization state.
// Also resets the service singleton so New() can be called again.
// Only available in test builds.
func Reset() {
	slush.Reset()
	if instance != nil {
		instance.mu.Lock()
		instance.encryptors = make(map[EncryptAlgo]Encryptor)
		instance.hashers = make(map[HashAlgo]Hasher)
		instance.maskers = make(map[MaskType]Masker)
		instance.codec = nil
		if instance.aperture != nil {
			instance.aperture.Close()
		}
		if instance.providers != nil {
			_ = instance.providers.Shutdown(context.Background())
		}
		instance.aperture = nil
		instance.providers = nil
		instance.mu.Unlock()
	}
	instance = nil
	once = sync.Once{}
}

// Unregister removes a service by type.
// Only available in test builds.
func Unregister[T any]() {
	slush.Unregister[T]()
}
