package main

import (
	"os"
)

type Config struct {
	DataDir       string
	RedisUser     string
	RedisPassword string
	RedisHost     string
	RedisPort     string
	RedisPoolSize int
	HostData      bool
	Port          string
}

// getEnvWithDefault retrieves variables from environment or uses
// default variables
func getEnvWithDefault(key string, i interface{}) interface{} {
	if val := os.Getenv(key); val == "" {
		return i
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
	DefaultHostData := true
	DefaultRedisPoolSize := 100

	return Config{
		DataDir:       getEnvWithDefault("DATA_DIR", DefaultDataDir).(string),
		RedisUser:     getEnvWithDefault("REDIS_USER", DefaultRedisUser).(string),
		RedisPassword: getEnvWithDefault("REDIS_PASSWORD", "").(string),
		RedisHost:     getEnvWithDefault("REDIS_HOST", DefaultRedisHost).(string),
		RedisPort:     getEnvWithDefault("REDIS_PORT", DefaultRedisPort).(string),
		RedisPoolSize: getEnvWithDefault("REDIS_POOL_SIZE", DefaultRedisPoolSize).(int),
		HostData:      getEnvWithDefault("HOST_DATA", DefaultHostData).(bool),
		Port:          getEnvWithDefault("PORT", DefaultAPIPort).(string),
	}
}
