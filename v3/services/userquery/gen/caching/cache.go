// Code generated by goa v3.0.9, DO NOT EDIT.
//
// Caching functions
//
// Command:
// $ goa gen gitlab.com/wiserskills/v3/services/userquery/design

package caching

import (
	"fmt"
	"strings"

	"github.com/dgraph-io/ristretto"
)

// Cacheable is implemented by cacheable entities/results
type Cacheable interface {
	CacheKey() string
}

// KeyBuilder is used to concatenate values and create a unique cache key
type KeyBuilder struct {
	values []string
}

// Cache represents a cache
type Cache struct {
	innerCache *ristretto.Cache
}

// NewCache creates a new cache
func NewCache() (*Cache, error) {

	c, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})

	if err != nil {
		return nil, err
	}

	return &Cache{innerCache: c}, nil
}

// NewKeyBuilder creates a new key builder
func NewKeyBuilder() *KeyBuilder {
	return &KeyBuilder{
		values: []string{},
	}
}

// Add adds a value to the key builder
func (c *KeyBuilder) Add(value interface{}) {
	str := strings.ReplaceAll(strings.ToLower(fmt.Sprintf("%v", value)), " ", "_")
	c.values = append(c.values, str)
}

// Get returns the key
func (c *KeyBuilder) Get() string {
	return strings.Join(c.values[:], ":")
}

// Set sets the cached value for the specified key
func (c *Cache) Set(key string, value interface{}) bool {

	return c.innerCache.Set(key, value, 1)
}

// SetCacheable sets a cacheable item
func (c *Cache) SetCacheable(item Cacheable) bool {

	return c.innerCache.Set(item.CacheKey(), item, 1)
}

// Get returns the cached value for the specified key
func (c *Cache) Get(key string) interface{} {

	value, found := c.innerCache.Get(key)

	if !found {
		return nil
	}

	return value
}

// Delete deletes the specified key from the cache
func (c *Cache) Delete(key string) {

	c.innerCache.Del(key)
}

// DeleteCacheable deletes the specified item from the cache
func (c *Cache) DeleteCacheable(item Cacheable) {

	c.innerCache.Del(item.CacheKey())
}

// Clear clears the cache
func (c *Cache) Clear() {

	c.innerCache.Clear()
}