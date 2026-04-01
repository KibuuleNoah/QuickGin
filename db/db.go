package db

import (
	"fmt"
	"log"
	"os"

	_redis "github.com/go-redis/redis/v7"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var dbConn *sqlx.DB

// Init connects to PostgreSQL using environment variables.
func Init() {
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
	if err != nil {
		log.Fatal(err)
	}
}

// GetDB returns the sqlx database connection.
func GetDB() *sqlx.DB {
	return dbConn
}

// RedisClient holds the Redis connection.
var RedisClient *_redis.Client

// InitRedis connects to Redis.
func InitRedis(selectDB ...int) {
	var redisHost = os.Getenv("REDIS_HOST")
	var redisPassword = os.Getenv("REDIS_PASSWORD")

	redisDB := 0
	if len(selectDB) > 0 {
		redisDB = selectDB[0]
	}

	RedisClient = _redis.NewClient(&_redis.Options{
		Addr:     redisHost,
		Password: redisPassword,
		DB:       redisDB,
	})

	if err := RedisClient.Ping().Err(); err != nil {
		log.Fatalf("Critical Error: Could not connect to Redis: %v", err)
	}

}

// GetRedis returns the Redis client.
func GetRedis() *_redis.Client {
	return RedisClient
}
