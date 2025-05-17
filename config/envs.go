package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// EnvNotFoundError is an error returned when an environment variable is not found
type EnvNotFoundError struct {
	EnvName string
}

func (e EnvNotFoundError) Error() string {
	return fmt.Sprintf("environment variable %s is not set", e.EnvName)
}

// GetEnv retrieves an environment variable value and validates if it's declared
// Returns the value if found, or an error if not set
func GetEnv(key string) (string, error) {
	value, exists := os.LookupEnv(key)

	if !exists || strings.TrimSpace(value) == "" {
		return "", EnvNotFoundError{EnvName: key}
	}

	return value, nil
}

// MustGetEnv retrieves an environment variable value and panics if it's not set
func MustGetEnv(key string) string {
	value, err := GetEnv(key)
	if err != nil {
		panic(err)
	}
	return value
}

// GetEnvAs retrieves and converts an environment variable to the specified type
func GetEnvAs[T any](key string) (T, error) {
	var result T
	value, err := GetEnv(key)
	if err != nil {
		return result, err
	}

	switch any(result).(type) {
	case string:
		return any(value).(T), nil
	case int:
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return result, fmt.Errorf("failed to convert %s to int: %w", key, err)
		}
		return any(intValue).(T), nil
	case int64:
		int64Value, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return result, fmt.Errorf("failed to convert %s to int64: %w", key, err)
		}
		return any(int64Value).(T), nil
	case float64:
		float64Value, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return result, fmt.Errorf("failed to convert %s to float64: %w", key, err)
		}
		return any(float64Value).(T), nil
	case bool:
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return result, fmt.Errorf("failed to convert %s to bool: %w", key, err)
		}
		return any(boolValue).(T), nil
	default:
		return result, fmt.Errorf("unsupported type for environment variable %s", key)
	}
}

// MustGetEnvAs retrieves and converts an environment variable to the specified type, panics if conversion fails
func MustGetEnvAs[T any](key string) T {
	value, err := GetEnvAs[T](key)
	if err != nil {
		panic(err)
	}
	return value
}
