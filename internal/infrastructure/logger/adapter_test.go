package logger

import (
	"errors"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/Soyuen/go-redis-chat-server/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestNewZapLogger_Success(t *testing.T) {
	cfg := config.AppConfig{
		IsDebug:      true,
		LogToFile:    false,
		LogFilePath:  "",
		LogErrorPath: "",
	}
	logger, err := NewZapLogger(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, logger)

	// Call several log methods to ensure they do not panic
	logger.Infow("info test")
	logger.Debugw("debug test")
	logger.Warnw("warn test")
	logger.Errorw("error test")
	// Fatalw calls os.Exit(1), so it's usually not called during testing.

	err = logger.Sync()
	if err != nil && !errors.Is(err, syscall.EBADF) {
		t.Errorf("unexpected Sync error: %v", err)
	}
}

func TestNewZapLogger_FailCreateDir(t *testing.T) {
	// Set an invalid path to cause mkdir to fail
	cfg := config.AppConfig{
		IsDebug:      false,
		LogToFile:    true,
		LogFilePath:  "/root/forbidden/log.json",
		LogErrorPath: "/root/forbidden/error.log",
	}

	logger, err := NewZapLogger(cfg)
	assert.Error(t, err)
	assert.Nil(t, logger)
}

func TestNewZapLogger_LogToFile(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "app.log")
	errorFile := filepath.Join(tmpDir, "error.log")

	cfg := config.AppConfig{
		IsDebug:      false,
		LogToFile:    true,
		LogFilePath:  logFile,
		LogErrorPath: errorFile,
	}

	logger, err := NewZapLogger(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, logger)

	logger.Infow("info log test")
	logger.Errorw("error log test")
	err = logger.Sync()
	assert.NoError(t, err)
}
