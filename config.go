package gopiapi

import (
    "os"
    "errors"
)

type Config struct {
    DataDir string
    RedisUser string
    RedisPassword string
    RedisHost string
    RedisPort string
    RedisPoolSize int
    Port string
}

func getEnvWithDefault(key string, d string) string {
    if val := os.Getenv(key); val == "" {
        return d
    } else {
        return val
    }
}

// GetConfigFromEnv retrieves API configuration details from environment.
// Provides defaults for all, including RedisPassword. Returns Config struct
func GetConfigFromEnv() Config {
    DefaultDataDir := "/mnt/external/data/"
    DefaultRedisUser := "redis"
    DefaultRedisPort := "6379"
    DefaultRedisHost := "localhost"
    DefaultAPIPort := "12345"
    DefaultRedisPoolSize := 100

    return Config{
        DataDir: getEnvWithDefault("DATA_DIR", DefaultDataDir),
        RedisUser: getEnvWithDefault("REDIS_USER", DefaultRedisUser),
        RedisPassword: getEnvWithDefault("REDIS_PASSWORD", ""),
        RedisHost: getEnvWithDefault("REDIS_HOST", DefaultRedisHost),
        RedisPort: getEnvWithDefault("REDIS_PORT", DefaultRedisPort),
        RedisPoolSize: getEnvWithDefault("REDIS_POOL_SIZE", DefaultRedisPoolSize),
        Port: getEnvWithDefault("PORT", DefaultAPIPort),
    }
}