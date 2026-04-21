package db

import (
	"fmt"
	"log"
	"os"

	"github.com/KibuuleNoah/QuickGin/models/cache"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

var dbConn *sqlx.DB
var appCache cache.Cache

// Init connects to PostgreSQL using environment variables.
func InitAppDB() {
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

	// dbConn, err = sqlx.Connect("sqlite", "file:sqlite3.db")
	if err != nil {
		log.Fatalln(err)
	}
}

// AppDB returns the sqlx database connection.
func AppDB() *sqlx.DB {
	return dbConn
}

func InitAppCache() {
	appCache = cache.NewMemoryCache()

}

func AppCache() cache.Cache {
	return appCache
}

// // GetRedis returns the Redis client.
// func GetRedis() *_redis.Client {
// 	return RedisClient
// }
