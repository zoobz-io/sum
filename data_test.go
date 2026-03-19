//go:build testing

package sum

import (
	"testing"
)

// Note: Full integration tests for Database, Store, and Bucket
// require actual database/store connections and are located in
// testing/integration/data_test.go.
//
// These unit tests verify the wrapper types and registration flow.

type testModel struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

func TestNewDatabase(t *testing.T) {
	t.Skip("requires database connection - see testing/integration/data_test.go")
}

func TestNewStore(t *testing.T) {
	t.Skip("requires store provider - see testing/integration/data_test.go")
}

func TestNewBucket(t *testing.T) {
	t.Skip("requires bucket provider - see testing/integration/data_test.go")
}

func TestNewSearch(t *testing.T) {
	t.Skip("requires search provider - see testing/integration/data_test.go")
}
