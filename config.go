package main

import (
	"os"
	"strconv"
)

type ServerConfig struct {
	ListenAddress string
}

type Config struct {
	Server       ServerConfig
	TestMode     bool
	SecretKey    string
	DatabaseURI  string
	DatabaseName string
}

var Conf *Config

func NewConfig() {
	Conf = &Config{
		Server: ServerConfig{
			ListenAddress: getEnv("LISTEN_ADDRESS", "localhost:8080"),
		},
		TestMode:     getEnvAsBool("TEST_MODE", false),
		SecretKey:    getEnv("SECRET_KEY", "very secret key"),
		DatabaseURI:  getEnv("DATABASE_URI", "mongodb://localhost:27017/"),
		DatabaseName: getEnv("DATABASE_NAME", "db"),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}