# Integration Tests

Integration tests verify sum components working together.

## Running

```bash
make test-integration
```

## Writing Integration Tests

- Use the `integration` build tag
- Test realistic scenarios with multiple components
- May use external resources (databases, services)
- Longer timeouts acceptable

Example:

```go
//go:build integration

package integration

import (
    "testing"
    sumtest "github.com/zoobz-io/sum/testing"
)

func TestServiceIntegration(t *testing.T) {
    ctx := sumtest.TestContext(t)
    // Integration test code
}
```
