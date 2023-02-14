package golog

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
	"time"
)

type Level = zapcore.Level

const (
	DebugLevel = zapcore.DebugLevel
	InfoLevel  = zapcore.InfoLevel
	WarnLevel  = zapcore.WarnLevel
	ErrorLevel = zapcore.ErrorLevel
)

func StringLevel(level string) Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	}

	return zapcore.InfoLevel
}

// LoggerInterface hide something you can implement new one
type LoggerInterface interface {
	SetOutputFile(logPath, fileName string) LoggerInterface
	SetFileRotate(fileMaxAge, fileRotation time.Duration) LoggerInterface
	SetLevel(level Level) LoggerInterface
	SetCallerShort(short bool) LoggerInterface
	SetName(name string) LoggerInterface
	SetIsOutputStdout(isOutputStdout bool) LoggerInterface
	SetCallerSkip(skip int) LoggerInterface
	SetOutputJson(json bool) LoggerInterface

	GetOutputFile() (logPath, fileName string)
	GetFileRotate() (fileMaxAge, fileRotation time.Duration)
	GetLevel() (level Level)
	GetCallerShort() (short bool)
	GetName() (name string)
	GetIsOutputStdout() (isOutputStdout bool)
	GetCallerSkip() (skip int)
	GetOutputJson() bool

	// InitLogger init logger should call this when change config
	InitLogger()
	// Sync terminal the logger should call this to flush
	Sync() error

	Panicf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Debugf(template string, args ...interface{})

	Panic(args ...interface{})
	Fatal(args ...interface{})
	Error(args ...interface{})
	Warn(args ...interface{})
	Info(args ...interface{})
	Debug(args ...interface{})

	PanicWithFields(fields map[string]interface{}, template string, args ...interface{})
	FatalWithFields(fields map[string]interface{}, template string, args ...interface{})
	ErrorWithFields(fields map[string]interface{}, template string, args ...interface{})
	WarnWithFields(fields map[string]interface{}, template string, args ...interface{})
	InfoWithFields(fields map[string]interface{}, template string, args ...interface{})
	DebugWithFields(fields map[string]interface{}, template string, args ...interface{})

	PanicContext(ctx context.Context, template string, args ...interface{})
	FatalContext(ctx context.Context, template string, args ...interface{})
	ErrorContext(ctx context.Context, template string, args ...interface{})
	WarnContext(ctx context.Context, template string, args ...interface{})
	InfoContext(ctx context.Context, template string, args ...interface{})
	DebugContext(ctx context.Context, template string, args ...interface{})

	PanicContextWithFields(ctx context.Context, fields map[string]interface{}, template string, args ...interface{})
	FatalContextWithFields(ctx context.Context, fields map[string]interface{}, template string, args ...interface{})
	ErrorContextWithFields(ctx context.Context, fields map[string]interface{}, template string, args ...interface{})
	WarnContextWithFields(ctx context.Context, fields map[string]interface{}, template string, args ...interface{})
	InfoContextWithFields(ctx context.Context, fields map[string]interface{}, template string, args ...interface{})
	DebugContextWithFields(ctx context.Context, fields map[string]interface{}, template string, args ...interface{})

	// AddFieldFunc filter deal the fields
	AddFieldFunc(func(context.Context, map[string]interface{}))

	GetZapLogger() *zap.Logger
	GetZapSugaredLogger() *zap.SugaredLogger
}
