//go:build testing

package integration

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/zoobz-io/astql/postgres"
	"github.com/zoobz-io/grub"
	"github.com/zoobz-io/sum"
	sumtest "github.com/zoobz-io/sum/testing"
)

// testModel represents a simple entity for testing database operations.
type testModel struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

// testDBStore embeds sum.Database for testing custom store patterns.
type testDBStore struct {
	*sum.Database[testModel]
}

func TestDatabaseIntegration(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL not set - skipping database integration test")
	}

	ctx := sumtest.TestContext(t)
	_ = ctx

	sum.Reset()
	t.Cleanup(sum.Reset)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	_ = sum.Start()

	// Create test table
	_, err = sqlxDB.Exec(`
		CREATE TABLE IF NOT EXISTS test_models (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL
		)
	`)
	if err != nil {
		t.Fatalf("failed to create test table: %v", err)
	}
	t.Cleanup(func() {
		sqlxDB.Exec(`DROP TABLE IF EXISTS test_models`)
	})

	database := sum.NewDatabase[testModel](sqlxDB, "test_models", postgres.New())

	if database.Database == nil {
		t.Error("expected non-nil embedded grub.Database")
	}

	// Verify embedding pattern works
	store := &testDBStore{Database: database}
	if store.Database == nil {
		t.Error("expected non-nil database in embedded store")
	}
}

// mockStoreProvider is a test implementation of grub.StoreProvider.
type mockStoreProvider struct{}

func (m mockStoreProvider) Get(ctx context.Context, key string) ([]byte, error) {
	return nil, nil
}

func (m mockStoreProvider) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return nil
}

func (m mockStoreProvider) Delete(ctx context.Context, key string) error {
	return nil
}

func (m mockStoreProvider) Exists(ctx context.Context, key string) (bool, error) {
	return false, nil
}

func (m mockStoreProvider) List(ctx context.Context, prefix string, limit int) ([]string, error) {
	return nil, nil
}

func (m mockStoreProvider) GetBatch(ctx context.Context, keys []string) (map[string][]byte, error) {
	return nil, nil
}

func (m mockStoreProvider) SetBatch(ctx context.Context, items map[string][]byte, ttl time.Duration) error {
	return nil
}

// testKVStore embeds sum.Store for testing custom store patterns.
type testKVStore struct {
	*sum.Store[testModel]
}

func TestStoreIntegration(t *testing.T) {
	// This test uses a mock provider. For actual KV store testing,
	// set TEST_STORE_PROVIDER and implement real provider initialization.
	if os.Getenv("TEST_STORE_PROVIDER") == "" {
		t.Skip("TEST_STORE_PROVIDER not set - skipping store integration test")
	}

	ctx := sumtest.TestContext(t)
	_ = ctx

	sum.Reset()
	t.Cleanup(sum.Reset)

	_ = sum.Start()

	store, err := sum.NewStore[testModel](mockStoreProvider{}, "test-store")
	if err != nil {
		t.Fatalf("NewStore failed: %v", err)
	}

	if store == nil {
		t.Fatal("expected non-nil store")
	}

	if store.Store == nil {
		t.Error("expected non-nil embedded grub.Store")
	}

	// Verify embedding pattern works
	kvStore := &testKVStore{Store: store}
	if kvStore.Store == nil {
		t.Error("expected non-nil store in embedded store")
	}
}

// mockBucketProvider is a test implementation of grub.BucketProvider.
type mockBucketProvider struct{}

func (m mockBucketProvider) Get(ctx context.Context, key string) ([]byte, *grub.ObjectInfo, error) {
	return nil, nil, nil
}

func (m mockBucketProvider) Put(ctx context.Context, key string, data []byte, info *grub.ObjectInfo) error {
	return nil
}

func (m mockBucketProvider) Delete(ctx context.Context, key string) error {
	return nil
}

func (m mockBucketProvider) Exists(ctx context.Context, key string) (bool, error) {
	return false, nil
}

func (m mockBucketProvider) List(ctx context.Context, prefix string, limit int) ([]grub.ObjectInfo, error) {
	return nil, nil
}

// testBlobStore embeds sum.Bucket for testing custom store patterns.
type testBlobStore struct {
	*sum.Bucket[testModel]
}

func TestBucketIntegration(t *testing.T) {
	// This test uses a mock provider. For actual blob storage testing,
	// set TEST_BUCKET_PROVIDER and implement real provider initialization.
	if os.Getenv("TEST_BUCKET_PROVIDER") == "" {
		t.Skip("TEST_BUCKET_PROVIDER not set - skipping bucket integration test")
	}

	ctx := sumtest.TestContext(t)
	_ = ctx

	sum.Reset()
	t.Cleanup(sum.Reset)

	_ = sum.Start()

	bucket, err := sum.NewBucket[testModel](mockBucketProvider{}, "test-bucket")
	if err != nil {
		t.Fatalf("NewBucket failed: %v", err)
	}

	if bucket == nil {
		t.Fatal("expected non-nil bucket")
	}

	if bucket.Bucket == nil {
		t.Error("expected non-nil embedded grub.Bucket")
	}

	// Verify embedding pattern works
	blobStore := &testBlobStore{Bucket: bucket}
	if blobStore.Bucket == nil {
		t.Error("expected non-nil bucket in embedded store")
	}
}
