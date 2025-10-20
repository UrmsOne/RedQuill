package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv   string
	HTTPPort int
	MongoURI string
	DBName   string
    JWTSecret string
    JWTTTLMin int
}

func Load() Config {
	_ = godotenv.Load()

	port := 8080
	if v := os.Getenv("PORT"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			port = p
		} else {
			log.Printf("invalid PORT, fallback to %d: %v", port, err)
		}
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

    return Config{
		AppEnv:   getenv("APP_ENV", "development"),
		HTTPPort: port,
		MongoURI: mongoURI,
        DBName:   getenv("MONGO_DB", "redquill"),
        JWTSecret: getenv("JWT_SECRET", "dev-secret-change-me"),
        JWTTTLMin: atoi(getenv("JWT_TTL_MIN", "120"), 120),
	}
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func atoi(s string, def int) int {
    if i, err := strconv.Atoi(s); err == nil {
        return i
    }
    return def
}


