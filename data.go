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
// Panics if catalog registration fails (duplicate table name is a programmer error).
func NewDatabase[M any](db *sqlx.DB, table string, renderer astql.Renderer) *Database[M] {
	gdb := grub.NewDatabase[M](db, table, renderer)
	if err := svc().catalog.RegisterDatabase("db://"+table, gdb.Atomic()); err != nil {
		panic("sum: NewDatabase: " + err.Error())
	}
	return &Database[M]{Database: gdb}
}

// Store wraps grub.Store and registers with scio on creation.
// Embed this type in your store structs to add custom methods.
type Store[M any] struct {
	*grub.Store[M]
}

// NewStore creates a Store[M] and registers it with the scio catalog.
// Requires sum.New() to have been called first.
// Panics if catalog registration fails (duplicate store name is a programmer error).
func NewStore[M any](provider grub.StoreProvider, name string) *Store[M] {
	store := grub.NewStore[M](provider)
	if err := svc().catalog.RegisterStore("kv://"+name, store.Atomic()); err != nil {
		panic("sum: NewStore: " + err.Error())
	}
	return &Store[M]{Store: store}
}

// Bucket wraps grub.Bucket and registers with scio on creation.
// Embed this type in your store structs to add custom methods.
type Bucket[M any] struct {
	*grub.Bucket[M]
}

// NewBucket creates a Bucket[M] and registers it with the scio catalog.
// Requires sum.New() to have been called first.
// Panics if catalog registration fails (duplicate bucket name is a programmer error).
func NewBucket[M any](provider grub.BucketProvider, name string) *Bucket[M] {
	bucket := grub.NewBucket[M](provider)
	if err := svc().catalog.RegisterBucket("bcs://"+name, bucket.Atomic()); err != nil {
		panic("sum: NewBucket: " + err.Error())
	}
	return &Bucket[M]{Bucket: bucket}
}

// Search wraps grub.Search and registers with scio on creation.
// Embed this type in your store structs to add custom query methods.
type Search[M any] struct {
	*grub.Search[M]
}

// NewSearch creates a Search[M] and registers it with the scio catalog.
// Requires sum.New() to have been called first.
// Panics if catalog registration fails (duplicate index name is a programmer error).
func NewSearch[M any](provider grub.SearchProvider, index string) *Search[M] {
	search := grub.NewSearch[M](provider, index)
	if err := svc().catalog.RegisterSearch("srch://"+index, search.Atomic()); err != nil {
		panic("sum: NewSearch: " + err.Error())
	}
	return &Search[M]{Search: search}
}
