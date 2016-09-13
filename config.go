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
	GatewayPort       = ":" + getHostValueOrDefault("GATEWAY_PORT_27017_TCP_PORT", getHostValueOrDefault("PORT", "80"))
	GatewayHost       = getHostValueOrDefault("GATEWAY_PORT_27017_TCP_ADDR", getHostValueOrDefault("IP_ADDRESS", ""))
	StartTesting bool = getHostValueOrDefault("TEST", "") != ""
)
