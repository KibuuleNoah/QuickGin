package cache

import (
	"encoding/json"
	"sync"
	"time"
)

type memoryCache struct {
	data sync.Map
}

func NewMemoryCache() Cache {
	mc := &memoryCache{}
	go mc.startSweeper() // Clean up expired keys automatically
	return mc
}

func (m *memoryCache) Get(key string) (CacheItem, bool) {
	val, ok := m.data.Load(key)
	if !ok {
		return CacheItem{}, false
	}

	item := val.(CacheItem)
	if time.Now().After(item.Expiration) {
		m.data.Delete(key)
		return CacheItem{}, false
	}

	// Return the JSON string so it mimics the Redis return type
	return item, true
}

func (m *memoryCache) Set(key string, value interface{}, exp time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	m.data.Store(key, CacheItem{
		Value:      jsonData,
		Expiration: time.Now().Add(exp),
	})
	return nil
}

func (m *memoryCache) Delete(key string) error {
	m.data.Delete(key)
	return nil
}

// startSweeper mimics Redis's active expiration cleanup
func (m *memoryCache) startSweeper() {
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		m.data.Range(func(key, value interface{}) bool {
			item := value.(CacheItem)
			if time.Now().After(item.Expiration) {
				m.data.Delete(key)
			}
			return true
		})
	}
}
