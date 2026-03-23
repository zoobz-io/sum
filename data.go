package sum

import (
	"github.com/jmoiron/sqlx"
	"github.com/zoobz-io/astql"
	"github.com/zoobz-io/grub"
)

// Database wraps grub.Database and registers with scio on creation.
// Embed this type in your store structs to add custom query methods.
type Database[M any] struct {
	*grub.Database[M]
}

// NewDatabase creates a Database[M] and registers it with the scio catalog.
// Requires sum.New() to have been called first.
func NewDatabase[M any](db *sqlx.DB, table string, renderer astql.Renderer) (*Database[M], error) {
	gdb, err := grub.NewDatabase[M](db, table, renderer)
	if err != nil {
		return nil, err
	}
	if err := svc().catalog.RegisterDatabase("db://"+table, gdb.Atomic()); err != nil {
		return nil, err
	}
	return &Database[M]{Database: gdb}, nil
}

// Store wraps grub.Store and registers with scio on creation.
// Embed this type in your store structs to add custom methods.
type Store[M any] struct {
	*grub.Store[M]
}

// NewStore creates a Store[M] and registers it with the scio catalog.
// Requires sum.New() to have been called first.
func NewStore[M any](provider grub.StoreProvider, name string) (*Store[M], error) {
	store := grub.NewStore[M](provider)
	if err := svc().catalog.RegisterStore("kv://"+name, store.Atomic()); err != nil {
		return nil, err
	}
	return &Store[M]{Store: store}, nil
}

// Bucket wraps grub.Bucket and registers with scio on creation.
// Embed this type in your store structs to add custom methods.
type Bucket[M any] struct {
	*grub.Bucket[M]
}

// NewBucket creates a Bucket[M] and registers it with the scio catalog.
// Requires sum.New() to have been called first.
func NewBucket[M any](provider grub.BucketProvider, name string) (*Bucket[M], error) {
	bucket := grub.NewBucket[M](provider)
	if err := svc().catalog.RegisterBucket("bcs://"+name, bucket.Atomic()); err != nil {
		return nil, err
	}
	return &Bucket[M]{Bucket: bucket}, nil
}

// Search wraps grub.Search and registers with scio on creation.
// Embed this type in your store structs to add custom query methods.
type Search[M any] struct {
	*grub.Search[M]
}

// NewSearch creates a Search[M] and registers it with the scio catalog.
// Requires sum.New() to have been called first.
func NewSearch[M any](provider grub.SearchProvider, index string) (*Search[M], error) {
	search, err := grub.NewSearch[M](provider, index)
	if err != nil {
		return nil, err
	}
	if err := svc().catalog.RegisterSearch("srch://"+index, search.Atomic()); err != nil {
		return nil, err
	}
	return &Search[M]{Search: search}, nil
}
