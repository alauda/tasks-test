package main

import (
	"os"
)

func getHostValueOrDefault(key string, defaultVal string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultVal
	}
	return value
}

var (
	//Server
	GatewayPort = ":" + getHostValueOrDefault("PORT", "80")
	GatewayHost   = getHostValueOrDefault("IP_ADDRESS", "")
)
