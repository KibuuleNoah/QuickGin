package db

import (
	"QuickGin/config"
	"QuickGin/models/cache"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

var dbConn *sqlx.DB
var appCache cache.Cache

// Init connects to PostgreSQL using environment variables.
func InitAppDB() error {
	// sslMode := "disable"
	// if  == "TRUE" {
	// 	sslMode = "require"
	// }
	//
	// dbinfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
	// 	os.Getenv("DB_HOST"),
	// 	os.Getenv("DB_PORT"),
	// 	os.Getenv("DB_USER"),
	// 	os.Getenv("DB_PASS"),
	// 	os.Getenv("DB_NAME"),
	// 	sslMode,
	// )
	log.Println(config.Get().DSN())
	var err error
	dbConn, err = sqlx.Connect("postgres", config.Get().DSN())
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
