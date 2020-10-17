package logger

import (
	"github.com/AleksK1NG/api-mc/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)

const (
	envLogOutput = "LOG_OUTPUT"
)

type Logger struct {
	*zap.Logger
}

// Return new app logger instance
func NewLogger(cfg *config.Config) (*Logger, error) {
	c := zap.Config{
		OutputPaths:       []string{getOutput()},
		Level:             zap.NewAtomicLevelAt(getLevel(cfg.Logger.Level)),
		Development:       cfg.Logger.Development,
		DisableCaller:     cfg.Logger.DisableCaller,
		DisableStacktrace: cfg.Logger.DisableStacktrace,
		Encoding:          cfg.Logger.Encoding,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "MESSAGE",
			CallerKey:    "CALLER",
			LevelKey:     "LEVEL",
			TimeKey:      "TIME",
			NameKey:      "NAME_KEY",
			EncodeLevel:  getEncodeLevel(cfg),
			EncodeTime:   zapcore.RFC3339TimeEncoder,
			EncodeCaller: zapcore.FullCallerEncoder,
		},
	}

	l, err := c.Build()
	if err != nil {
		return nil, err
	}
	defer l.Sync()

	return &Logger{l}, nil
}

// Return response error and log actual error
func (l *Logger) ErrorWithLog(err error, responseError error) error {
	l.Error(err.Error())
	return responseError
}

func getEncodeLevel(cfg *config.Config) zapcore.LevelEncoder {
	if cfg.Logger.Encoding == "console" {
		return zapcore.CapitalColorLevelEncoder
	}
	return zapcore.CapitalLevelEncoder
}

func getLevel(level string) zapcore.Level {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "error":
		return zap.ErrorLevel
	default:
		return zap.InfoLevel
	}
}

func getOutput() string {
	output := strings.TrimSpace(os.Getenv(envLogOutput))
	if output == "" {
		return "stdout"
	}
	return output
}
