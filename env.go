package main

import "os"

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

var envSharedKey = getEnv("WHEREIS_SHARED_KEY", "")
var envListenAddr = getEnv("WHEREIS_LISTEN_ADDRESS", ":3000")
var envDbPath = getEnv("WHEREIS_DB_PATH", "./whereis.db3")
var envLocationIqKey = getEnv("WHEREIS_LIQ_KEY", "")
