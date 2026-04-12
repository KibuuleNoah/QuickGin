package cache

import (
	"sync"
	"time"
)

type memoryCache struct {
	data sync.Map
}

func NewMemoryCache() Cache {
	return &memoryCache{}
}

func (m *memoryCache) Get(key string) (string, error) {
	m.data.Load(key)
	return "", nil
}

func (m *memoryCache) Set(key string, value interface{}, exp time.Duration) error {
	m.data.Store(key, value)
	return nil
}

func (m *memoryCache) Delete(key string) error {
	m.data.Delete(key)
	return nil
}
