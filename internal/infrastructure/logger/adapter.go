package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Soyuen/go-redis-chat-server/internal/config"
	"github.com/Soyuen/go-redis-chat-server/pkg/loggeriface"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	sugar *zap.SugaredLogger
}

func NewZapLogger(cfg config.AppConfig) (loggeriface.Logger, error) {
	var cores []zapcore.Core

	encoderCfg := zap.NewProductionEncoderConfig()
	if cfg.IsDebug {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	}

	encoder := zapcore.NewJSONEncoder(encoderCfg)

	// Console output (always include if IsDebug)
	if cfg.IsDebug {
		consoleCore := zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderCfg),
			zapcore.AddSync(os.Stdout),
			zapcore.DebugLevel,
		)
		cores = append(cores, consoleCore)
	}

	if cfg.LogToFile {
		if err := os.MkdirAll(filepath.Dir(cfg.LogFilePath), os.ModePerm); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}
		infoFileWriter, err := os.OpenFile(cfg.LogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}

		errorFileWriter, err := os.OpenFile(cfg.LogErrorPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open error log file: %w", err)
		}

		infoCore := zapcore.NewCore(
			encoder,
			zapcore.AddSync(infoFileWriter),
			zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl < zapcore.ErrorLevel
			}),
		)

		errorCore := zapcore.NewCore(
			encoder,
			zapcore.AddSync(errorFileWriter),
			zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl >= zapcore.ErrorLevel
			}),
		)

		cores = append(cores, infoCore, errorCore)
	}

	combinedCore := zapcore.NewTee(cores...)

	logger := zap.New(combinedCore, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return &ZapLogger{sugar: logger.Sugar()}, nil
}

func (z *ZapLogger) Infow(msg string, keysAndValues ...interface{}) {
	z.sugar.Infow(msg, keysAndValues...)
}

func (z *ZapLogger) Debugw(msg string, keysAndValues ...interface{}) {
	z.sugar.Debugw(msg, keysAndValues...)
}

func (z *ZapLogger) Warnw(msg string, keysAndValues ...interface{}) {
	z.sugar.Warnw(msg, keysAndValues...)
}

func (z *ZapLogger) Errorw(msg string, keysAndValues ...interface{}) {
	z.sugar.Errorw(msg, keysAndValues...)
}

func (z *ZapLogger) Fatalw(msg string, keysAndValues ...interface{}) {
	z.sugar.Fatalw(msg, keysAndValues...)
}

func (z *ZapLogger) Sync() error {
	return z.sugar.Sync()
}
