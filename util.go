package main

import (
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func envOrFlag[T any](envKey string, flagValue T) T {
	if value, exist := os.LookupEnv(envKey); exist {
		switch any(flagValue).(type) {
		case string:
			return any(value).(T)
		case bool:
			boolValue, err := strconv.ParseBool(value)
			if err != nil {
				panic(err)
			}
			return any(boolValue).(T)
		case time.Duration:
			durationValue, err := time.ParseDuration(value)
			if err != nil {
				panic(err)
			}
			return any(durationValue).(T)
		default:
			panic("unsupported conversion type")
		}
	}

	return flagValue
}

func matchPattern(str string, pattern string) (bool, error) {
	if str == pattern {
		return true, nil
	}

	r, err := regexp.Compile(pattern)
	if err != nil {
		return false, err
	}

	match := r.MatchString(str)
	if match {
		if strings.Contains(str, pattern) {
			match = false
		}
	}

	return match, nil
}
