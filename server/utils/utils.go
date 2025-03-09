package utils

import (
	"log"
	"os"
)


func ReqEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("Missing required environment variable: %s", k)
	}

	return v
}
