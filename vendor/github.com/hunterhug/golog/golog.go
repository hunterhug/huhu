package golog

import (
	"context"
	rotateLogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type logger struct {
	name         string
	zapLogger    *zap.Logger
	sugarLog     *zap.SugaredLogger
	level        Level
	short        bool
	json         bool
	addFieldFunc func(context.Context, map[string]interface{})

	logPath  string
	fileName string

	fileMaxAge, fileRotation time.Duration

	isOutputStdout bool
	skip           int
}

var _log = New()

func init() {
	// skip is 2, we wrap 2 layer
	_log.SetCallerSkip(2)
	_log.InitLogger()
}

// Logger default log which output to console
// if you want to log to file you must New() and set something then call InitLogger()
func Logger() LoggerInterface {
	return _log
}

// New can new a logger interface, you can config it by it's method
func New() LoggerInterface {
	l := new(logger)
	l.level = InfoLevel
	l.short = false
	l.json = false
	return l
}

// InitLogger after config you must call this method
func (l *logger) InitLogger() {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.LevelKey = "l"
	encoderConfig.FunctionKey = "func"
	encoderConfig.CallerKey = "caller"
	encoderConfig.TimeKey = "t"
	encoderConfig.MessageKey = "msg"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	if l.short {
		encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	} else {
		encoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	}

	encoderConfig.LineEnding = zapcore.DefaultLineEnding
	var zConfig zapcore.Encoder

	if !l.json {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		zConfig = zapcore.NewConsoleEncoder(encoderConfig)

	} else {
		zConfig = zapcore.NewJSONEncoder(encoderConfig)
	}

	var outCore zapcore.Core

	if l.logPath != "" {
		cores := make([]zapcore.Core, 0)
		debugFileName := filepath.Join(l.logPath, "access.log")
		infoFileName := filepath.Join(l.logPath, "info.log")
		warnFileName := filepath.Join(l.logPath, "warn.log")
		errFileName := filepath.Join(l.logPath, "error.log")

		if l.fileName != "" {
			debugFileName = filepath.Join(l.logPath, l.fileName) + "_debug.log"
			infoFileName = filepath.Join(l.logPath, l.fileName) + "_info.log"
			warnFileName = filepath.Join(l.logPath, l.fileName) + "_warn.log"
			errFileName = filepath.Join(l.logPath, l.fileName) + "_err.log"
		}

		if l.level <= DebugLevel {
			debugLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl >= DebugLevel
			})

			debugWriter := getWriter(false, debugFileName, l.fileMaxAge, l.fileRotation)
			debugWriteSync := zapcore.AddSync(debugWriter)
			core := zapcore.NewCore(
				zConfig,
				debugWriteSync,
				debugLevel,
			)

			cores = append(cores, core)
		}

		if l.level <= InfoLevel {
			infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl >= InfoLevel
			})

			infoWriter := getWriter(false, infoFileName, l.fileMaxAge, l.fileRotation)
			infoWriteSync := zapcore.AddSync(infoWriter)
			core := zapcore.NewCore(
				zConfig,
				infoWriteSync,
				infoLevel,
			)
			cores = append(cores, core)
		}

		if l.level <= WarnLevel {
			warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl >= WarnLevel
			})

			warnWriter := getWriter(false, warnFileName, l.fileMaxAge, l.fileRotation)
			warnWriteSync := zapcore.AddSync(warnWriter)
			core := zapcore.NewCore(
				zConfig,
				warnWriteSync,
				warnLevel,
			)
			cores = append(cores, core)
		}

		if l.level <= ErrorLevel {
			errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl >= ErrorLevel
			})

			errorWriter := getWriter(false, errFileName, l.fileMaxAge, l.fileRotation)
			errorWriteSync := zapcore.AddSync(errorWriter)
			core := zapcore.NewCore(
				zConfig,
				errorWriteSync,
				errorLevel,
			)
			cores = append(cores, core)

		}

		if l.isOutputStdout {
			stdOutLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl >= l.level
			})
			core := zapcore.NewCore(
				zConfig,
				zapcore.AddSync(os.Stdout),
				stdOutLevel,
			)
			cores = append(cores, core)
		}
		outCore = zapcore.NewTee(cores...)

	} else {
		writeSync := zapcore.AddSync(os.Stdout)
		outCore = zapcore.NewCore(
			zConfig,
			writeSync,
			l.level,
		)
	}

	op1 := zap.AddCaller()

	// we wrap 1 layer
	op2 := zap.AddCallerSkip(1)

	if l.skip > 0 {
		op2 = zap.AddCallerSkip(l.skip)
	}

	zapLogger := zap.New(outCore, op1, op2)
	if l.name != "" {
		zapLogger = zapLogger.Named(l.name)
	}
	sugarLogger := zapLogger.Sugar()
	l.zapLogger = zapLogger
	l.sugarLog = sugarLogger
}

func InitLogger() {
	_log.InitLogger()
}

func getWriter(isOutputStdout bool, filename string, maxAge, rotation time.Duration) io.Writer {
	if maxAge <= 0 {
		maxAge = 30 * 24 * time.Hour
	}

	if rotation <= 0 {
		rotation = 24 * time.Hour
	}

	format := "%Y%m%d%H%M.log"
	if rotation == 24*time.Hour {
		format = "%Y%m%d.log"
	}

	hook, err := rotateLogs.New(
		filename+"."+format,
		rotateLogs.WithLinkName(filename),
		rotateLogs.WithMaxAge(maxAge),
		rotateLogs.WithRotationTime(rotation),
	)

	if err != nil {
		panic(err)
	}

	if isOutputStdout {
		return io.MultiWriter(os.Stdout, hook)
	}

	return hook
}

func (l *logger) Sync() error {
	return l.sugarLog.Sync()
}

func Sync() error {
	return _log.Sync()
}

func SetName(name string) LoggerInterface {
	return _log.SetName(name)
}

func (l *logger) SetName(name string) LoggerInterface {
	l.name = name
	return l
}

func GetName() (name string) {
	return _log.GetName()
}

func (l *logger) GetName() (name string) {
	return l.name
}

func SetCallerCallerSkip(skip int) LoggerInterface {
	return _log.SetCallerSkip(skip)
}

func (l *logger) SetCallerSkip(skip int) LoggerInterface {
	l.skip = skip
	return l
}

func GetCallerSkip() (skip int) {
	return _log.GetCallerSkip()
}

func (l *logger) GetCallerSkip() (skip int) {
	return l.skip
}

func SetIsOutputStdout(isOutputStdout bool) LoggerInterface {
	return _log.SetIsOutputStdout(isOutputStdout)
}

func (l *logger) SetIsOutputStdout(isOutputStdout bool) LoggerInterface {
	l.isOutputStdout = isOutputStdout
	return l
}

func GetIsOutputStdout() (isOutputStdout bool) {
	return _log.GetIsOutputStdout()
}

func (l *logger) GetIsOutputStdout() (isOutputStdout bool) {
	return l.isOutputStdout
}

func SetFileRotate(fileMaxAge, fileRotation time.Duration) LoggerInterface {
	return _log.SetFileRotate(fileMaxAge, fileRotation)
}

func (l *logger) SetFileRotate(fileMaxAge, fileRotation time.Duration) LoggerInterface {
	l.fileMaxAge = fileMaxAge
	l.fileRotation = fileRotation
	return l
}

func GetFileRotate() (fileMaxAge, fileRotation time.Duration) {
	return _log.GetFileRotate()
}

func (l *logger) GetFileRotate() (fileMaxAge, fileRotation time.Duration) {
	return l.fileMaxAge, l.fileRotation
}

func SetCallerShort(short bool) LoggerInterface {
	return _log.SetCallerShort(short)
}

func (l *logger) SetCallerShort(short bool) LoggerInterface {
	l.short = short
	return l
}

func GetCallerShort() (short bool) {
	return _log.GetCallerShort()
}

func (l *logger) GetCallerShort() (short bool) {
	return l.short
}

func SetOutputJson(json bool) LoggerInterface {
	return _log.SetOutputJson(json)
}

func (l *logger) SetOutputJson(json bool) LoggerInterface {
	l.json = json
	return l
}

func GetOutputJson() (json bool) {
	return _log.GetOutputJson()
}

func (l *logger) GetOutputJson() (json bool) {
	return l.json
}

func SetLevel(level Level) LoggerInterface {
	return _log.SetLevel(level)
}

func (l *logger) SetLevel(level Level) LoggerInterface {
	l.level = level
	return l
}

func GetLevel() (level Level) {
	return _log.GetLevel()
}

func (l *logger) GetLevel() (level Level) {
	return l.level
}

func (l *logger) SetOutputFile(logPath, fileName string) LoggerInterface {
	l.logPath = logPath
	l.fileName = fileName
	if l.fileMaxAge == 0 {
		l.fileMaxAge = 30 * 24 * time.Hour
		l.fileRotation = 24 * time.Hour
	}
	return l
}

func SetOutputFile(logPath, fileName string) LoggerInterface {
	return _log.SetOutputFile(logPath, fileName)
}

func (l *logger) GetOutputFile() (logPath, fileName string) {
	return l.logPath, l.fileName
}

func GetOutputFile() (logPath, fileName string) {
	return _log.GetOutputFile()
}

func (l *logger) Fatalf(template string, args ...interface{}) {
	l.sugarLog.Fatalf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	_log.Fatalf(template, args...)
}

func (l *logger) Fatal(args ...interface{}) {
	l.sugarLog.Fatal(args...)
}

func Fatal(args ...interface{}) {
	_log.Fatal(args...)
}

func (l *logger) Panicf(template string, args ...interface{}) {
	l.sugarLog.With().Panicf(template, args...)
}

func Panicf(template string, args ...interface{}) {
	_log.Panicf(template, args...)
}

func (l *logger) Panic(args ...interface{}) {
	l.sugarLog.Panic(args...)
}

func Panic(args ...interface{}) {
	_log.Panic(args...)
}

func (l *logger) Errorf(template string, args ...interface{}) {
	l.sugarLog.Errorf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	_log.Errorf(template, args...)
}

func (l *logger) Error(args ...interface{}) {
	l.sugarLog.Error(args...)
}

func Error(args ...interface{}) {
	_log.Error(args...)
}

func (l *logger) Warnf(template string, args ...interface{}) {
	l.sugarLog.Warnf(template, args...)
}

func Warnf(template string, args ...interface{}) {
	_log.Warnf(template, args...)
}

func (l *logger) Warn(args ...interface{}) {
	l.sugarLog.Warn(args...)
}

func Warn(args ...interface{}) {
	_log.Warn(args...)
}

func (l *logger) Infof(template string, args ...interface{}) {
	l.sugarLog.Infof(template, args...)
}

func Infof(template string, args ...interface{}) {
	_log.Infof(template, args...)
}

func (l *logger) Info(args ...interface{}) {
	l.sugarLog.Info(args...)
}

func Info(args ...interface{}) {
	_log.Info(args...)
}

func (l *logger) Debugf(template string, args ...interface{}) {
	l.sugarLog.Debugf(template, args...)
}

func Debugf(template string, args ...interface{}) {
	_log.Debugf(template, args...)
}

func (l *logger) Debug(args ...interface{}) {
	l.sugarLog.Debug(args...)
}

func Debug(args ...interface{}) {
	_log.Debug(args...)
}

func with(fields map[string]interface{}) []interface{} {
	i := make([]interface{}, 0, 2*len(fields))
	keys := make([]string, 0, len(fields))
	for k, _ := range fields {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	for _, v := range keys {
		i = append(i, v, fields[v])
	}
	return i
}

func (l *logger) DebugWithFields(fields map[string]interface{}, template string, args ...interface{}) {
	if len(args) == 0 {
		l.sugarLog.With(with(fields)...).Debug(template)
		return
	}
	l.sugarLog.With(with(fields)...).Debugf(template, args...)
}

func DebugWithFields(fields map[string]interface{}, template string, args ...interface{}) {
	_log.DebugWithFields(fields, template, args...)
}

func (l *logger) InfoWithFields(fields map[string]interface{}, template string, args ...interface{}) {
	if len(args) == 0 {
		l.sugarLog.With(with(fields)...).Info(template)
		return
	}
	l.sugarLog.With(with(fields)...).Infof(template, args...)
}

func InfoWithFields(fields map[string]interface{}, template string, args ...interface{}) {
	_log.InfoWithFields(fields, template, args...)
}

func (l *logger) WarnWithFields(fields map[string]interface{}, template string, args ...interface{}) {
	if len(args) == 0 {
		l.sugarLog.With(with(fields)...).Warn(template)
		return
	}
	l.sugarLog.With(with(fields)...).Warnf(template, args...)
}

func WarnWithFields(fields map[string]interface{}, template string, args ...interface{}) {
	_log.WarnWithFields(fields, template, args...)
}

func (l *logger) ErrorWithFields(fields map[string]interface{}, template string, args ...interface{}) {
	if len(args) == 0 {
		l.sugarLog.With(with(fields)...).Error(template)
		return
	}
	l.sugarLog.With(with(fields)...).Errorf(template, args...)
}

func ErrorWithFields(fields map[string]interface{}, template string, args ...interface{}) {
	_log.ErrorWithFields(fields, template, args...)
}

func (l *logger) FatalWithFields(fields map[string]interface{}, template string, args ...interface{}) {
	if len(args) == 0 {
		l.sugarLog.With(with(fields)...).Fatal(template)
		return
	}
	l.sugarLog.With(with(fields)...).Fatalf(template, args...)
}

func FatalWithFields(fields map[string]interface{}, template string, args ...interface{}) {
	_log.FatalWithFields(fields, template, args...)
}

func (l *logger) PanicWithFields(fields map[string]interface{}, template string, args ...interface{}) {
	if len(args) == 0 {
		l.sugarLog.With(with(fields)...).Panic(template)
		return
	}
	l.sugarLog.With(with(fields)...).Panicf(template, args...)
}

func PanicWithFields(fields map[string]interface{}, template string, args ...interface{}) {
	_log.PanicWithFields(fields, template, args...)
}

func AddFieldFunc(f func(context.Context, map[string]interface{})) {
	_log.AddFieldFunc(f)
}

func (l *logger) AddFieldFunc(f func(context.Context, map[string]interface{})) {
	l.addFieldFunc = f
	return
}

func (l *logger) addField(ctx context.Context, fields map[string]interface{}) {
	//fields["service.log.name"] = l.name
	//fields["service.log.time"] = time.Now().String()

	if l.addFieldFunc != nil {
		l.addFieldFunc(ctx, fields)
	}
}

func (l *logger) DebugContextWithFields(ctx context.Context, fields map[string]interface{}, template string, args ...interface{}) {
	l.addField(ctx, fields)
	if len(args) == 0 {
		l.sugarLog.With(with(fields)...).Debug(template)
		return
	}
	l.sugarLog.With(with(fields)...).Debugf(template, args...)
}

func DebugContextWithFields(ctx context.Context, fields map[string]interface{}, template string, args ...interface{}) {
	_log.DebugContextWithFields(ctx, fields, template, args...)
}

func (l *logger) InfoContextWithFields(ctx context.Context, fields map[string]interface{}, template string, args ...interface{}) {
	l.addField(ctx, fields)
	if len(args) == 0 {
		l.sugarLog.With(with(fields)...).Info(template)
		return
	}
	l.sugarLog.With(with(fields)...).Infof(template, args...)
}

func InfoContextWithFields(ctx context.Context, fields map[string]interface{}, template string, args ...interface{}) {
	_log.InfoContextWithFields(ctx, fields, template, args...)
}

func (l *logger) WarnContextWithFields(ctx context.Context, fields map[string]interface{}, template string, args ...interface{}) {
	l.addField(ctx, fields)
	if len(args) == 0 {
		l.sugarLog.With(with(fields)...).Warn(template)
		return
	}
	l.sugarLog.With(with(fields)...).Warnf(template, args...)
}

func WarnContextWithFields(ctx context.Context, fields map[string]interface{}, template string, args ...interface{}) {
	_log.WarnContextWithFields(ctx, fields, template, args...)
}

func (l *logger) ErrorContextWithFields(ctx context.Context, fields map[string]interface{}, template string, args ...interface{}) {
	l.addField(ctx, fields)
	if len(args) == 0 {
		l.sugarLog.With(with(fields)...).Error(template)
		return
	}
	l.sugarLog.With(with(fields)...).Errorf(template, args...)
}

func ErrorContextWithFields(ctx context.Context, fields map[string]interface{}, template string, args ...interface{}) {
	_log.ErrorContextWithFields(ctx, fields, template, args...)
}

func (l *logger) FatalContextWithFields(ctx context.Context, fields map[string]interface{}, template string, args ...interface{}) {
	l.addField(ctx, fields)
	if len(args) == 0 {
		l.sugarLog.With(with(fields)...).Fatal(template)
		return
	}
	l.sugarLog.With(with(fields)...).Fatalf(template, args...)
}

func FatalContextWithFields(ctx context.Context, fields map[string]interface{}, template string, args ...interface{}) {
	_log.FatalContextWithFields(ctx, fields, template, args...)
}

func (l *logger) PanicContextWithFields(ctx context.Context, fields map[string]interface{}, template string, args ...interface{}) {
	l.addField(ctx, fields)
	if len(args) == 0 {
		l.sugarLog.With(with(fields)...).Panic(template)
		return
	}
	l.sugarLog.With(with(fields)...).Panicf(template, args...)
}

func PanicContextWithFields(ctx context.Context, fields map[string]interface{}, template string, args ...interface{}) {
	_log.PanicContextWithFields(ctx, fields, template, args...)
}

func (l *logger) DebugContext(ctx context.Context, template string, args ...interface{}) {
	fields := make(map[string]interface{})
	l.addField(ctx, fields)
	if len(args) == 0 {
		l.sugarLog.With(with(fields)...).Debug(template)
		return
	}
	l.sugarLog.With(with(fields)...).Debugf(template, args...)
}

func DebugContext(ctx context.Context, template string, args ...interface{}) {
	_log.DebugContext(ctx, template, args...)
}

func (l *logger) InfoContext(ctx context.Context, template string, args ...interface{}) {
	fields := make(map[string]interface{})
	l.addField(ctx, fields)
	if len(args) == 0 {
		l.sugarLog.With(with(fields)...).Info(template)
		return
	}
	l.sugarLog.With(with(fields)...).Infof(template, args...)
}

func InfoContext(ctx context.Context, template string, args ...interface{}) {
	_log.InfoContext(ctx, template, args...)
}

func (l *logger) WarnContext(ctx context.Context, template string, args ...interface{}) {
	fields := make(map[string]interface{})
	l.addField(ctx, fields)
	if len(args) == 0 {
		l.sugarLog.With(with(fields)...).Warn(template)
		return
	}
	l.sugarLog.With(with(fields)...).Warnf(template, args...)
}

func WarnContext(ctx context.Context, template string, args ...interface{}) {
	_log.WarnContext(ctx, template, args...)
}

func (l *logger) ErrorContext(ctx context.Context, template string, args ...interface{}) {
	fields := make(map[string]interface{})
	l.addField(ctx, fields)
	if len(args) == 0 {
		l.sugarLog.With(with(fields)...).Error(template)
		return
	}
	l.sugarLog.With(with(fields)...).Errorf(template, args...)
}

func ErrorContext(ctx context.Context, template string, args ...interface{}) {
	_log.ErrorContext(ctx, template, args...)
}

func (l *logger) FatalContext(ctx context.Context, template string, args ...interface{}) {
	fields := make(map[string]interface{})
	l.addField(ctx, fields)
	if len(args) == 0 {
		l.sugarLog.With(with(fields)...).Fatal(template)
		return
	}
	l.sugarLog.With(with(fields)...).Fatalf(template, args...)
}

func FatalContext(ctx context.Context, template string, args ...interface{}) {
	_log.FatalContext(ctx, template, args...)
}

func (l *logger) PanicContext(ctx context.Context, template string, args ...interface{}) {
	fields := make(map[string]interface{})
	l.addField(ctx, fields)
	if len(args) == 0 {
		l.sugarLog.With(with(fields)...).Panic(template)
		return
	}
	l.sugarLog.With(with(fields)...).Panicf(template, args...)
}

func PanicContext(ctx context.Context, template string, args ...interface{}) {
	_log.PanicContext(ctx, template, args...)
}

func (l *logger) GetZapLogger() *zap.Logger {
	return l.zapLogger
}

func GetZapLogger() *zap.Logger {
	return _log.GetZapLogger()
}

func GetZapSugaredLogger() *zap.SugaredLogger {
	return _log.GetZapSugaredLogger()
}

func (l *logger) GetZapSugaredLogger() *zap.SugaredLogger {
	return l.sugarLog
}
