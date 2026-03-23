# Contributing to sum

Thank you for your interest in contributing to sum.

## Prerequisites

- Go 1.24 or later
- golangci-lint (installed via `make install-tools`)

## Development Setup

```bash
# Clone the repository
git clone https://github.com/zoobz-io/sum.git
cd sum

# Install development tools
make install-tools

# Install git hooks
make install-hooks
```

## Development Workflow

```bash
# Run tests
make test

# Run linter
make lint

# Quick validation (tests + lint)
make check

# Full CI simulation
make ci
```

Run `make help` for all available commands.

## Pull Request Process

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Ensure `make check` passes
5. Submit a pull request

## Code Standards

- All code must pass linting (`make lint`)
- All tests must pass (`make test`)
- New code requires tests (80% patch coverage target)
- Follow existing code patterns and conventions
