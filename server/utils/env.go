package utils

import (
	"log"
	"os"
	"strconv"
)

func GetNumericEnv(key string) int {
	env := os.Getenv(key)
	value, err := strconv.Atoi(env)
	if err != nil {
		log.Fatalf("Environment variable %s is supposed to be an integer. Got: %v", key, env)
	}

	return value
}
