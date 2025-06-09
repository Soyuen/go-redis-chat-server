package config

import "strconv"

type AppConfig struct {
	Port         string
	IsDebug      bool
	LogToFile    bool
	LogFilePath  string
	LogErrorPath string
}

func LoadEnvConfig() AppConfig {
	port := GetEnv("PORT", "8080")

	isDebugStr := GetEnv("APP_DEBUG", "false")

	isDebug, err := strconv.ParseBool(isDebugStr)
	if err != nil {
		isDebug = false
	}

	logToFileStr := GetEnv("LOG_TO_FILE", "true")
	logToFile, err := strconv.ParseBool(logToFileStr)
	if err != nil {
		logToFile = false
	}

	logFilePath := GetEnv("LOG_FILE_PATH", "logs/app.log")
	logErrorPath := GetEnv("LOG_ERROR_PATH", "logs/app.log")

	return AppConfig{
		Port:         port,
		IsDebug:      isDebug,
		LogToFile:    logToFile,
		LogFilePath:  logFilePath,
		LogErrorPath: logErrorPath,
	}
}
