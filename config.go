package main

import (
	"os"
	"strconv"
)

type Config struct {
	DataDir       string
	RedisPassword string
	RedisHost     string
	RedisPort     string
	RedisPoolSize int
	HostData      bool
	LogLevel      string
	Port          string
}

type getEnvArg struct {
	key     string
	Default interface{}
	T       string
}

func convertToInt(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}

	return i
}

func convertToType(typeString string, value string) interface{} {
	if typeString == "int" {
		return convertToInt(value)
	} else if typeString == "bool" {
		if v, err := strconv.ParseBool(value); err != nil {
			panic(err)
		} else {
			return v
		}
	}

	return value
}

// getEnvWithDefault retrieves variables from environment or uses
// default variables
func getEnvWithDefault(Arg getEnvArg) interface{} {
	if val := os.Getenv(Arg.key); val == "" {
		return Arg.Default
	} else {
		return convertToType(Arg.T, val)
	}
}

// GetConfigFromEnv retrieves API configuration details from environment.
// Provides defaults for all, including RedisPassword. Returns Config struct
func GetConfigFromEnv() Config {
	return Config{
		DataDir: getEnvWithDefault(getEnvArg{
			key:     "DATA_DIR",
			Default: "/mnt/external/data/",
			T:       "string",
		}).(string),
		RedisPassword: getEnvWithDefault(getEnvArg{
			key:     "REDIS_PASSWORD",
			Default: "",
			T:       "string",
		}).(string),
		RedisHost: getEnvWithDefault(getEnvArg{
			key:     "REDIS_HOST",
			Default: "127.0.0.1",
			T:       "string",
		}).(string),
		RedisPort: getEnvWithDefault(getEnvArg{
			key:     "REDIS_PORT",
			Default: "6379",
			T:       "string",
		}).(string),
		RedisPoolSize: getEnvWithDefault(getEnvArg{
			key:     "REDIS_POOL_SIZE",
			Default: 100,
			T:       "int",
		}).(int),
		HostData: getEnvWithDefault(getEnvArg{
			key:     "HOST_DATA",
			Default: true,
			T:       "bool",
		}).(bool),
		LogLevel: getEnvWithDefault(getEnvArg{
			key:     "LOGLEVEL",
			Default: "INFO",
			T:       "string",
		}).(string),
		Port: getEnvWithDefault(getEnvArg{
			key:     "PORT",
			Default: "12345",
			T:       "string",
		}).(string),
	}
}
