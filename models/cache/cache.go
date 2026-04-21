package cache

import (
	"encoding/json"
	"time"
)

type Cache interface {
	Get(key string) (CacheItem, bool)
	Set(key string, value interface{}, expiration time.Duration) error
	Delete(key string) error
}

type CacheItem struct {
	Value      []byte    `db:"value"`
	Expiration time.Time `db:"expires_at"`
}

// IsExpired checks if the item has passed its TTL
func (c *CacheItem) IsExpired() bool {
	return !c.Expiration.IsZero() && time.Now().After(c.Expiration)
}

// String returns the value as a string.
func (c *CacheItem) String() (string, error) {
	var s string
	err := json.Unmarshal(c.Value, &s)
	return s, err
}

// Int returns the value as an integer.
func (c *CacheItem) Int() (int, error) {
	var i int
	err := json.Unmarshal(c.Value, &i)
	return i, err
}

// Bool returns the value as a boolean.
func (c *CacheItem) Bool() (bool, error) {
	var b bool
	err := json.Unmarshal(c.Value, &b)
	return b, err
}

// Map returns the value as a map[string]interface{}.
func (c *CacheItem) Map() (map[string]interface{}, error) {
	var m map[string]interface{}
	err := json.Unmarshal(c.Value, &m)
	return m, err
}

// Interface (Any) returns the raw interface{} representation.
func (c *CacheItem) Interface() (interface{}, error) {
	var val interface{}
	err := json.Unmarshal(c.Value, &val)
	return val, err
}

// UnmarshalTo allows you to pass a custom struct pointer to fill
func (c *CacheItem) UnmarshalTo(v interface{}) error {
	return json.Unmarshal(c.Value, v)
}
