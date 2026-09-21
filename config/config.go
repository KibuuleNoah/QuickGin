package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config is the application config
type Config struct {
	Env  string // DEV | PROD
	Port string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	DBMaxConns int

	RedisAddr     string
	RedisPassword string
	RedisDB       int

	JWTSecret     string
	JWTExpiry     time.Duration
	RefreshExpiry time.Duration

	SMTPHost string
	SMTPPort string
	SMTPUser string
	SMTPPass string

	RateLimitRPS   int
	RateLimitBurst int
}

var cfg *Config

// Load loads the config
func Load() *Config {
	if cfg != nil {
		return cfg
	}

	env := getEnv("APP_ENV", "DEV")

	// only load .env file in non-production; prod should inject real env vars
	if env != "PROD" {
		if err := godotenv.Load(); err != nil {
			log.Println("no .env file found, relying on system env")
		}
	}

	c := &Config{
		Env:  env,
		Port: getEnv("PORT", "8080"),

		DBHost:     mustGetEnv("DB_HOST"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     mustGetEnv("DB_USER"),
		DBPassword: mustGetEnv("DB_PASSWORD"),
		DBName:     mustGetEnv("DB_NAME"),
		DBSSLMode:  getEnv("DB_SSLMODE", "require"),
		DBMaxConns: getEnvInt("DB_MAX_CONNS", 20),

		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvInt("REDIS_DB", 0),

		JWTSecret:     mustGetEnv("JWT_SECRET"),
		JWTExpiry:     getEnvDuration("JWT_EXPIRY", 15*time.Minute),
		RefreshExpiry: getEnvDuration("REFRESH_EXPIRY", 7*24*time.Hour),

		SMTPHost: getEnv("SMTP_HOST", ""),
		SMTPPort: getEnv("SMTP_PORT", "587"),
		SMTPUser: getEnv("SMTP_USER", ""),
		SMTPPass: getEnv("SMTP_PASS", ""),

		RateLimitRPS:   getEnvInt("RATE_LIMIT_RPS", 10),
		RateLimitBurst: getEnvInt("RATE_LIMIT_BURST", 20),
	}

	if c.Env == "PROD" && c.JWTSecret == "changeme" {
		log.Fatal("refusing to start: default JWT_SECRET in production")
	}

	cfg = c
	return cfg
}

func Get() *Config {
	if cfg == nil {
		log.Fatal("config not loaded — call config.Load() first")
	}
	return cfg
}

func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode)
}

func (c *Config) IsProd() bool {
	return c.Env == "PROD"
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func mustGetEnv(key string) string {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		log.Fatalf("missing required env var: %s", key)
	}
	return v
}

func getEnvInt(key string, fallback int) int {
	v, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		log.Fatalf("invalid int for %s: %v", key, err)
	}
	return i
}

// getEnvDuration returns a time.Duration from an env var, or fallback if not set
func getEnvDuration(key string, fallback time.Duration) time.Duration {
	v, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		log.Fatalf("invalid duration for %s: %v", key, err)
	}
	return d
}
