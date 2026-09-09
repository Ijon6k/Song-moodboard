package config

import (
	"os"
)

type Config struct {
	Port              string
	DatabaseURL       string
	RedisAddr         string
	RedisPassword     string
	SeaweedEndpoint   string
	SeaweedAccessKey  string
	SeaweedSecretKey  string
	SeaweedBucket     string
	SeaweedPublicHost string
	JWTSecret         string
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultVal
}

func LoadConfig() *Config {
	dbHost := getEnv("DB_HOST", "postgres")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "songmoodboard")
	dbPass := getEnv("DB_PASSWORD", "songmoodboard_secret")
	dbName := getEnv("DB_NAME", "songmoodboard_db")
	dbSSL := getEnv("DB_SSLMODE", "disable")

	dbURL := getEnv("DATABASE_URL", "")
	if dbURL == "" {
		dbURL = "postgres://" + dbUser + ":" + dbPass + "@" + dbHost + ":" + dbPort + "/" + dbName + "?sslmode=" + dbSSL
	}

	return &Config{
		Port:              getEnv("PORT", "8080"),
		DatabaseURL:       dbURL,
		RedisAddr:         getEnv("REDIS_ADDR", "redis:6379"),
		RedisPassword:     getEnv("REDIS_PASSWORD", ""),
		SeaweedEndpoint:   getEnv("SEAWEED_ENDPOINT", "seaweedfs:8333"),
		SeaweedAccessKey:  getEnv("SEAWEED_ACCESS_KEY", "anykey"),
		SeaweedSecretKey:  getEnv("SEAWEED_SECRET_KEY", "anysecret"),
		SeaweedBucket:     getEnv("SEAWEED_BUCKET", "songmoodboard"),
		SeaweedPublicHost: getEnv("SEAWEED_PUBLIC_HOST", "/api/storage"), // routed via nginx or seaweedfs
		JWTSecret:         getEnv("JWT_SECRET", "super-secret-twilight-haze-key-2026"),
	}
}
