package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
)

func printWarning(msg string, args ...any) {
	fmt.Print("\033[33m")
	slog.Warn(msg, args...)
	fmt.Print("\033[0m")
}

func getString(envVarName, defaultValue string) string {
	envVar, exists := os.LookupEnv(envVarName)
	if !exists {
		printWarning(
			"Using default value for environment variable",
			slog.String("envVarName", envVarName),
			slog.String("defaultValue", defaultValue),
		)
		return defaultValue
	}
	return envVar
}

func getInt(envVarName string, defaultValue int) int {
	envVar, exists := os.LookupEnv(envVarName)
	if !exists {
		printWarning(
			"Using default value for environment variable",
			slog.String("envVarName", envVarName),
			slog.Int("defaultValue", defaultValue),
		)
		return defaultValue
	}
	intValue, err := strconv.ParseInt(envVar, 10, 16)
	if err != nil {
		printWarning(
			"Failed to parse environment variable",
			slog.String("envVarName", envVarName),
			slog.String("envVarValue", envVar),
			slog.Int("defaultValue", defaultValue),
		)
		return defaultValue
	}
	return int(intValue)
}


func getFloat32(envVarName string, defaultValue float32) float32 {
	envVar, exists := os.LookupEnv(envVarName)
	if !exists {
		printWarning(
			"Using default value for environment variable",
			slog.String("envVarName", envVarName),
			slog.Float64("defaultValue", float64(defaultValue)),
		)
		return defaultValue
	}
	floatValue, err := strconv.ParseFloat(envVar, 32)
	if err != nil {
		printWarning(
			"Failed to parse environment variable",
			slog.String("envVarName", envVarName),
			slog.String("envVarValue", envVar),
			slog.Float64("defaultValue", float64(defaultValue)),
		)
		return defaultValue
	}
	return float32(floatValue)
}
