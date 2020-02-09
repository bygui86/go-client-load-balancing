package utils

import (
	"os"
	"strconv"
)

// Check -
func Check(key string) (string, bool) {
	return os.LookupEnv(key)
}

// GetString -
func GetString(key, def string) string {
	value, exist := Check(key)
	if !exist {
		return def
	}
	return value
}

// GetBool -
func GetBool(key string, defaultValue bool) bool {
	value, err := strconv.ParseBool(os.Getenv(key))
	if err != nil {
		return defaultValue
	}
	return value
}

// GetInt -
func GetInt(key string, defaultValue int) int {
	return int(GetInt64(key, int64(defaultValue)))
}

// GetInt64 -
func GetInt64(key string, defaultValue int64) int64 {
	value, err := strconv.ParseInt(os.Getenv(key), 10, 64)
	if err != nil {
		return defaultValue
	}
	return value
}
