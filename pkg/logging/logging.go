package logging

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogWrapper struct {
	ZapLogger *zap.Logger
}
type LogModel struct {
	fields []zap.Field
}

var logger *LogWrapper

type LogConfig struct {
	Name     string
	Level    string
	RunMode  string
	Hostname string
	Encoding string
}

func Setup(config LogConfig) {
	var options []zap.Option
	logWriter := zapcore.AddSync(os.Stdout)

	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		NameKey:        "name",
		CallerKey:      "caller",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
	atomicLevel := zap.NewAtomicLevel()
	level, exist := loggerLevelMap[config.Level]
	if !exist {
		level = zapcore.DebugLevel
	}
	atomicLevel.SetLevel(level)

	var encoder zapcore.Encoder
	if config.Encoding == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}
	core := zapcore.NewCore(encoder, logWriter, atomicLevel)
	serviceNameField := zap.Fields(zap.String("serviceName", config.Name))
	hostNameField := zap.Fields(zap.String("hostname", config.Hostname))

	options = append(options, zap.AddCaller(), serviceNameField, hostNameField, zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))

	if config.Encoding == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	if config.RunMode != "release" {
		options = append(options, zap.Development())
	}
	logger = &LogWrapper{ZapLogger: zap.New(core, options...)}
}

var loggerLevelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

func GetLogger() *LogWrapper {
	return logger
}

func (l *LogWrapper) Debug(msg string, fields ...zap.Field) {
	fields = append(fields, zap.String("time", time.Now().UTC().Format(time.RFC3339)))
	logModel := LogModel{}
	logModel.fields = append(logModel.fields, fields...)
	l.ZapLogger.Debug(msg, logModel.fields...)
}

func (l *LogWrapper) Info(msg string, fields ...zap.Field) {
	fields = append(fields, zap.String("time", time.Now().UTC().Format(time.RFC3339)))
	logModel := LogModel{}
	logModel.fields = append(logModel.fields, fields...)
	l.ZapLogger.Info(msg, logModel.fields...)
}

func (l *LogWrapper) Warn(msg string, fields ...zap.Field) {
	fields = append(fields, zap.String("time", time.Now().UTC().Format(time.RFC3339)))
	logModel := LogModel{}
	logModel.fields = append(logModel.fields, fields...)
	l.ZapLogger.Warn(msg, logModel.fields...)
}

func (l *LogWrapper) Error(msg string, fields ...zap.Field) {
	fields = append(fields, zap.String("time", time.Now().UTC().Format(time.RFC3339)))
	logModel := LogModel{}
	logModel.fields = append(logModel.fields, fields...)
	l.ZapLogger.Error(msg, logModel.fields...)
}

func (l *LogWrapper) DPanic(msg string, fields ...zap.Field) {
	fields = append(fields, zap.String("time", time.Now().UTC().Format(time.RFC3339)))
	logModel := LogModel{}
	logModel.fields = append(logModel.fields, fields...)
	l.ZapLogger.DPanic(msg, logModel.fields...)
}
func (l *LogWrapper) Panic(msg string, fields ...zap.Field) {
	fields = append(fields, zap.String("time", time.Now().UTC().Format(time.RFC3339)))
	logModel := LogModel{}
	logModel.fields = append(logModel.fields, fields...)
	l.ZapLogger.Panic(msg, logModel.fields...)
}
func (l *LogWrapper) Fatal(msg string, fields ...zap.Field) {
	fields = append(fields, zap.String("time", time.Now().UTC().Format(time.RFC3339)))
	logModel := LogModel{}
	logModel.fields = append(logModel.fields, fields...)
	l.ZapLogger.Fatal(msg, logModel.fields...)
}
func (l *LogWrapper) Sync() error {
	return l.ZapLogger.Sync()
}
