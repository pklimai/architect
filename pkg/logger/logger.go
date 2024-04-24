package logger

import (
	"fmt"
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	global             *zap.SugaredLogger
	globalLevelEnabler = zap.NewAtomicLevelAt(zap.ErrorLevel)
)

func init() {
	SetLogger(New(globalLevelEnabler))
}

func New(levelEnabler zapcore.LevelEnabler, options ...zap.Option) *zap.SugaredLogger {
	return NewWithOutput(levelEnabler, os.Stdout, options...)
}

func NewWithOutput(
	levelEnabler zapcore.LevelEnabler,
	output io.Writer,
	opions ...zap.Option,
) *zap.SugaredLogger {
	if levelEnabler == nil {
		levelEnabler = globalLevelEnabler
	}

	cfg := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "ts",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		StacktraceKey:  "StacktraceKey",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	encoder := zapcore.NewJSONEncoder(cfg)

	core := zapcore.NewCore(encoder, zapcore.AddSync(output), levelEnabler)

	return zap.New(core, opions...).Sugar()
}

func Level() zapcore.Level {
	return globalLevelEnabler.Level()
}

func ParseLevel(levelString string) (zapcore.Level, error) {
	level, err := zapcore.ParseLevel(levelString)
	if err != nil {
		return zapcore.ErrorLevel, fmt.Errorf("zapcore.ParseLevel: %w", err)
	}

	return level, nil
}

func SetLevel(level zapcore.Level) {
	globalLevelEnabler.SetLevel(level)
}

func Logger() *zap.SugaredLogger {
	return global
}

func SetLogger(l *zap.SugaredLogger) {
	global = l
}

func Debug(args ...interface{}) {
	if global.Level().Enabled(zapcore.DebugLevel) {
		global.Debug(args...)
	}
}

func Debugf(format string, args ...interface{}) {
	if global.Level().Enabled(zapcore.DebugLevel) {
		global.Debugf(format, args...)
	}
}

func DebugKV(msg string, kvs ...interface{}) {
	if global.Level().Enabled(zapcore.DebugLevel) {
		global.Debugw(msg, kvs...)
	}
}

func Info(args ...interface{}) {
	if global.Level().Enabled(zapcore.InfoLevel) {
		global.Info(args...)
	}
}

func Infof(format string, args ...interface{}) {
	if global.Level().Enabled(zapcore.InfoLevel) {
		global.Infof(format, args...)
	}
}

func InfoKV(msg string, kvs ...interface{}) {
	if global.Level().Enabled(zapcore.InfoLevel) {
		global.Infow(msg, kvs...)
	}
}

func Warn(args ...interface{}) {
	if global.Level().Enabled(zapcore.WarnLevel) {
		global.Warn(args...)
	}
}

func Warnf(format string, args ...interface{}) {
	if global.Level().Enabled(zapcore.WarnLevel) {
		global.Warnf(format, args...)
	}
}

func WarnKV(msg string, kvs ...interface{}) {
	if global.Level().Enabled(zapcore.WarnLevel) {
		global.Warnw(msg, kvs...)
	}
}

func Error(args ...interface{}) {
	if global.Level().Enabled(zapcore.ErrorLevel) {
		global.Error(args...)
	}
}

func Errorf(format string, args ...interface{}) {
	if global.Level().Enabled(zapcore.ErrorLevel) {
		global.Errorf(format, args...)
	}
}

func ErrorKV(msg string, kvs ...interface{}) {
	if global.Level().Enabled(zapcore.ErrorLevel) {
		global.Errorw(msg, kvs...)
	}
}

func Fatal(args ...interface{}) {
	global.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	global.Fatalf(format, args...)
}

func FatalKV(msg string, kvs ...interface{}) {
	global.Fatalw(msg, kvs...)
}

func Panic(args ...interface{}) {
	global.Panic(args...)
}

func Panicf(format string, args ...interface{}) {
	global.Panicf(format, args...)
}

func PanicKV(msg string, kvs ...interface{}) {
	global.Panicw(msg, kvs...)
}
