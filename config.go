// Package sum provides an applications framework for Go.
package sum

import (
	"context"

	"github.com/zoobz-io/fig"
)

// Config loads configuration of type T via fig and registers it with the service locator.
// Pass nil for provider if secrets are not needed.
// Retrieve the configuration later with Use[T](ctx).
func Config[T any](ctx context.Context, k Key, provider fig.SecretProvider) error {
	var cfg T
	var opts []fig.SecretProvider
	if provider != nil {
		opts = append(opts, provider)
	}
	if err := fig.LoadContext(ctx, &cfg, opts...); err != nil {
		return err
	}
	Register[T](k, cfg)
	return nil
}
