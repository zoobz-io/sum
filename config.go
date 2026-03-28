// Package sum provides an applications framework for Go.
package sum

import (
	"context"
	"fmt"
	"reflect"

	"github.com/zoobz-io/fig"
)

// Config loads configuration of type T via fig and registers it with the service locator.
// Pass nil for provider if secrets are not needed.
// Retrieve the configuration later with Use[T](ctx).
// Errors are annotated with the type name for diagnostics.
func Config[T any](ctx context.Context, k Key, provider fig.SecretProvider) error {
	var cfg T
	var opts []fig.SecretProvider
	if provider != nil {
		opts = append(opts, provider)
	}
	if err := fig.LoadContext(ctx, &cfg, opts...); err != nil {
		return fmt.Errorf("config %s: %w", reflect.TypeOf(cfg).Name(), err)
	}
	Register[T](k, cfg)
	return nil
}

// ConfigAll runs configuration loaders sequentially and returns the first error.
// Each loader is typically a closure wrapping a Config[T] call:
//
//	sum.ConfigAll(
//	    func() error { return sum.Config[DBConfig](ctx, k, nil) },
//	    func() error { return sum.Config[RedisConfig](ctx, k, secrets) },
//	)
func ConfigAll(loaders ...func() error) error {
	for _, load := range loaders {
		if err := load(); err != nil {
			return err
		}
	}
	return nil
}
