package db

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"pajo/models/cache"

	_ "github.com/lib/pq"
)

var dbConn *sqlx.DB
var appCache cache.Cache

// Init connects to PostgreSQL using environment variables.
func InitAppDB() error {
	sslMode := "disable"
	if os.Getenv("SSL") == "TRUE" {
		sslMode = "require"
	}

	dbinfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
		sslMode,
	)

	var err error
	dbConn, err = sqlx.Connect("postgres", dbinfo)
	return err
}

// AppDB returns the sqlx database connection.
func AppDB() *sqlx.DB {
	return dbConn
}

type CacheType string

const (
	RedisCache    CacheType = "redis"
	MemCache      CacheType = "memory"
	PostgresCache CacheType = "postgres"
)

func InitAppCache(cacheType CacheType) error {

	switch cacheType {
	case MemCache:
		appCache = cache.NewMemoryCache()

	case PostgresCache:
		appCache = cache.NewPostgresCache(AppDB())

	case RedisCache:

		// appCache = cache.NewRedisCache()

	default:
		return fmt.Errorf("unsupported cache type")
	}
	return nil
}

func AppCache() cache.Cache {
	return appCache
}
