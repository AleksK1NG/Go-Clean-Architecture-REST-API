package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)

const (
	envLogLevel  = "LOG_LEVEL"
	envLogOutput = "LOG_OUTPUT"
)

type Logger struct {
	*zap.Logger
}

// Return new app logger instance
func NewLogger() (*Logger, error) {
	config := zap.Config{
		OutputPaths:       []string{getOutput()},
		Level:             zap.NewAtomicLevel(),
		Development:       true,
		DisableCaller:     false,
		DisableStacktrace: false,
		Encoding:          "console",
		EncoderConfig: zapcore.EncoderConfig{
			CallerKey:    "CALLER",
			LevelKey:     "LEVEL",
			TimeKey:      "TIME",
			NameKey:      "MESSAGE",
			EncodeLevel:  zapcore.CapitalColorLevelEncoder,
			EncodeTime:   zapcore.RFC3339TimeEncoder,
			EncodeCaller: zapcore.FullCallerEncoder,
		},
	}

	l, err := config.Build()
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

func getLevel() zapcore.Level {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(envLogLevel))) {
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
