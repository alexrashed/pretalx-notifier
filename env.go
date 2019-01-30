package main

import (
	"fmt"
	"os"
	"strconv"
)

// returns the string value of an environment variable or an error if not set or empty
func getEnvStr(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return value, fmt.Errorf(fmt.Sprintf("environment variable %s empty", key))
	}
	return value, nil
}

// returns the int value of an environment variable or an error if not set or empty
func getEnvInt(key string) (int, error) {
	str, err := getEnvStr(key)
	if err != nil {
		return 0, err
	}
	value, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}
	return value, nil
}

// returns the boole value of an environment variable or an error if not set or empty
func getEnvBool(key string) (bool, error) {
	str, err := getEnvStr(key)
	if err != nil {
		return false, err
	}
	value, err := strconv.ParseBool(str)
	if err != nil {
		return false, err
	}
	return value, nil
}