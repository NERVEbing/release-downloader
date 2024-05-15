package main

import (
	"os"
	"regexp"
	"strings"
	"time"
)

func getEnvOrDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvOrDefaultBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	if value == "true" || value == "1" {
		return true
	}

	return false
}

func getEnvOrDefaultDuration(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	d, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}

	return d
}

func matchWildcard(str string, pattern string) (bool, error) {
	if strings.Contains(str, pattern) {
		return true, nil
	}

	return regexp.MatchString(pattern, str)
}
