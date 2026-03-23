# sum

[![CI Status](https://github.com/zoobz-io/sum/workflows/CI/badge.svg)](https://github.com/zoobz-io/sum/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/zoobz-io/sum/graph/badge.svg?branch=main)](https://codecov.io/gh/zoobz-io/sum)
[![Go Report Card](https://goreportcard.com/badge/github.com/zoobz-io/sum)](https://goreportcard.com/report/github.com/zoobz-io/sum)
[![CodeQL](https://github.com/zoobz-io/sum/workflows/CodeQL/badge.svg)](https://github.com/zoobz-io/sum/security/code-scanning)
[![Go Reference](https://pkg.go.dev/badge/github.com/zoobz-io/sum.svg)](https://pkg.go.dev/github.com/zoobz-io/sum)
[![License](https://img.shields.io/github/license/zoobz-io/sum)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/zoobz-io/sum)](go.mod)
[![Release](https://img.shields.io/github/v/release/zoobz-io/sum)](https://github.com/zoobz-io/sum/releases)

**Wire once, run anywhere.** An application framework that unifies HTTP, data, configuration, and services into a single lifecycle.

## Compose and Run

```go
// Register services by contract type
k := sum.Start()
sum.Register[UserService](k, &userImpl{})
sum.Register[OrderService](k, &orderImpl{})
sum.Freeze(k)

// Retrieve anywhere by type
userSvc := sum.MustUse[UserService](ctx)
```

Services, configuration, and data stores—all wired through one registry, resolved by type.

## Install

```bash
go get github.com/zoobz-io/sum
```

Requires Go 1.24 or later.

## Quick Start

```go
package main

import (
    "context"
    "log"

    "github.com/zoobz-io/sum"
)

type Greeter interface {
    Greet(name string) string
}

type greeterImpl struct{}

func (g *greeterImpl) Greet(name string) string {
    return "Hello, " + name
}

func main() {
    // Initialize service and registry
    svc := sum.New(sum.ServiceConfig{Host: "localhost", Port: 8080})
    k := sum.Start()

    // Register services
    sum.Register[Greeter](k, &greeterImpl{})
    sum.Freeze(k)

    // Use services anywhere
    greeter := sum.MustUse[Greeter](context.Background())
    log.Println(greeter.Greet("World"))

    // Run with graceful shutdown
    if err := svc.Run(); err != nil {
        log.Fatal(err)
    }
}
```

## Capabilities

| Capability | Description | Documentation |
|------------|-------------|---------------|
| Service Registry | Type-safe service locator with guards | [Registry](https://pkg.go.dev/github.com/zoobz-io/sum#Register) |
| Lifecycle Management | Singleton service with graceful shutdown | [Service](https://pkg.go.dev/github.com/zoobz-io/sum#Service) |
| Configuration | Load and register config via fig | [Config](https://pkg.go.dev/github.com/zoobz-io/sum#Config) |
| Typed Events | Emit and listen with type-safe payloads | [Event](https://pkg.go.dev/github.com/zoobz-io/sum#Event) |
| Data Stores | Database, KV, and object storage helpers | [Database](https://pkg.go.dev/github.com/zoobz-io/sum#Database) |

## Why sum?

- **Type-safe service registry** — Register and retrieve services by contract type, not strings. Compile-time safety, zero casting.
- **Unified lifecycle** — One `Run()` handles startup, signal handling, and graceful shutdown.
- **Integrated data catalog** — Databases, KV stores, and buckets register automatically with the data catalog for observability.
- **Typed events** — Define events once with `Event[T]`, emit and listen with full type safety.
- **Minimal ceremony** — No annotations, no reflection magic, no code generation. Just Go.

## The Ecosystem

sum builds on the zoobz-io toolkit:

| Package | Purpose |
|---------|---------|
| [rocco](https://github.com/zoobz-io/rocco) | HTTP engine with OpenAPI |
| [slush](https://github.com/zoobz-io/slush) | Service registry core |
| [capitan](https://github.com/zoobz-io/capitan) | Event/signal system |
| [fig](https://github.com/zoobz-io/fig) | Configuration loading |
| [grub](https://github.com/zoobz-io/grub) | Database/KV/Object storage |
| [scio](https://github.com/zoobz-io/scio) | Data catalog |

## Documentation

- **Learn**: [Overview](docs/1.learn/1.overview.md) · [Quickstart](docs/1.learn/2.quickstart.md) · [Concepts](docs/1.learn/3.concepts.md) · [Architecture](docs/1.learn/4.architecture.md)
- **Guides**: [Testing](docs/2.guides/1.testing.md) · [Troubleshooting](docs/2.guides/2.troubleshooting.md) · [Service Registry](docs/2.guides/3.service-registry.md) · [Events](docs/2.guides/4.events.md) · [Data Stores](docs/2.guides/5.data-stores.md)
- **Reference**: [API](docs/4.reference/1.api.md) · [Types](docs/4.reference/2.types.md) · [pkg.go.dev](https://pkg.go.dev/github.com/zoobz-io/sum)

## Contributing

Contributions welcome—see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

MIT License - see [LICENSE](LICENSE)
