package main

import (
	"os"
	"strconv"
	"strings"
)

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

type Config struct {
	// TODO: implement this as enum
	role     string
	httpPort string
	nodes    []string
	delay    int
}

func (config *Config) init() {
	config.httpPort = getEnv("HTTP_PORT", "8080")
	config.role = getEnv("ROLE", "follower")
	config.delay, _ = strconv.Atoi(getEnv("DELAY", "0"))
	config.nodes = strings.Split(os.Getenv("FOLLOWERS"), ",")

	if len(config.nodes) == 0 {
		panic("provide follower list")
	}
}

var config = Config{"", "", []string{}, 0}
