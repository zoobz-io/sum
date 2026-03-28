package sum

import (
	"context"
	"fmt"
	"io"

	"github.com/zoobz-io/capitan"
	"github.com/zoobz-io/sentinel"
)

// Signals emitted during infrastructure lifecycle.
var (
	SignalConnected    = capitan.NewSignal("sum.connected", "Infrastructure client connected")
	SignalDisconnected = capitan.NewSignal("sum.disconnected", "Infrastructure client disconnected")
)

// Field keys for connector signals.
var (
	KeyConnectorName  = capitan.NewStringKey("connector")
	KeyConnectorType  = capitan.NewStringKey("type")
	KeyConnectorError = capitan.NewErrorKey("error")
)

// namedCloser pairs a name with an io.Closer for ordered shutdown.
type namedCloser struct {
	name   string
	closer io.Closer
}

// Connect creates an infrastructure client via factory, registers it with slush,
// emits a connected signal, and tracks io.Closer instances for automatic shutdown.
// The name identifies this connection in signals and shutdown logs.
// If the client implements io.Closer, it will be closed during Service.Shutdown
// in reverse connection order.
func Connect[T any](ctx context.Context, k Key, name string, factory func(context.Context) (T, error)) error {
	client, err := factory(ctx)
	if err != nil {
		return fmt.Errorf("connect %s: %w", name, err)
	}
	Register[T](k, client)
	typeName := name
	if meta, err := sentinel.TryInspect[T](); err == nil {
		typeName = meta.FQDN
	}
	capitan.Info(ctx, SignalConnected,
		KeyConnectorName.Field(name),
		KeyConnectorType.Field(typeName),
	)

	if closer, ok := any(client).(io.Closer); ok {
		s := svc()
		s.mu.Lock()
		s.closers = append(s.closers, namedCloser{name: name, closer: closer})
		s.mu.Unlock()
	}
	return nil
}
