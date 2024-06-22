package internal

import (
	"os"
	"strconv"
	"strings"
)

func GetAsString(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

func GetAsBool(name string, defaultValue bool) bool {
	valStr := GetAsString(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}
	return defaultValue
}

func GetAsInt(name string, defaultValue int) int {
	valueStr := GetAsString(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func GetAsSlice(name string, defaultValue []string, sep string) []string {
	valStr := GetAsString(name, "")

	if valStr == "" {
		return defaultValue
	}
	return strings.Split(valStr, sep)
}
