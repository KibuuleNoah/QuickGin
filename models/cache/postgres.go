package cache

import (
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

type PostgresCache struct {
	DB *sqlx.DB
}

func NewPostgresCache(DB *sqlx.DB) Cache {
	return &PostgresCache{
		DB: DB,
	}
}

func (p *PostgresCache) Get(key string) (CacheItem, bool) {
	var item CacheItem

	// Get automatically maps columns to struct fields via db tags
	// We check expiration directly in SQL so the app never sees "dead" data
	query := `SELECT value, expires_at FROM cache_items 
	          WHERE key = $1 AND (expires_at > NOW() OR expires_at IS NULL)`

	err := p.DB.Get(&item, query, key)
	if err != nil {
		return CacheItem{}, false // Returns false if key missing or expired
	}

	return item, true
}

func (p *PostgresCache) Set(key string, value interface{}, expiration time.Duration) error {
	valBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	var expiresAt interface{}
	if expiration > 0 {
		expiresAt = time.Now().Add(expiration)
	}
	// if expiration <= 0, expiresAt remains nil (NULL in DB)

	query := `INSERT INTO cache_items (key, value, expires_at) 
	          VALUES ($1, $2, $3)
	          ON CONFLICT (key) DO UPDATE SET value = $2, expires_at = $3`

	_, err = p.DB.Exec(query, key, valBytes, expiresAt)
	return err
}

func (p *PostgresCache) Delete(key string) error {
	_, err := p.DB.Exec("DELETE FROM cache_items WHERE key = $1", key)
	return err
}
