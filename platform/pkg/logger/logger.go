package logger

import (
	"context"
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Key string

const (
	traceIDKey Key = "trace_id"
	userIDKey  Key = "user_id"
)

var (
	// Глобавльный синглтон логгер
	globalLogger *WrapLogger
	// объект синхронизации, выполняющий код внутри единоразово.
	// Используется для реализации паттерна сингтон.
	initOnce sync.Once
	// Текущий уровень логгирования (из zapcore), который можно менять в runtime
	dynamicLevel zap.AtomicLevel
)

// WrapLogger - кастомная обертка над zap.Logger
type WrapLogger struct {
	zapLogger *zap.Logger
}

// Init - создание экземпляра логгера
func Init(levelStr string, asJSON bool) {
	initOnce.Do(
		func() {
			dynamicLevel = zap.NewAtomicLevelAt(parseLevel(levelStr))

			encoderCfg := buildProductionEncoderConfig()

			var encoder zapcore.Encoder
			if asJSON {
				encoder = zapcore.NewJSONEncoder(encoderCfg)
			} else {
				encoder = zapcore.NewConsoleEncoder(encoderCfg)
			}

			core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), dynamicLevel)

			zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2))

			globalLogger = &WrapLogger{zapLogger: zapLogger}
		},
	)
}

// Logger - получение логгера для использования в вашем коде
func Logger() *WrapLogger {
	return globalLogger
}

// SetNopLogger устанавливает глобальный логгер в no-op режим.
// Идеально для юнит-тестов.
func SetNopLogger() {
	globalLogger = &WrapLogger{
		zapLogger: zap.NewNop(),
	}
}

// SetLevel динамически меняет уровень логирования
func SetLevel(levelStr string) {
	if dynamicLevel == (zap.AtomicLevel{}) {
		return
	}

	dynamicLevel.SetLevel(parseLevel(levelStr))
}

// With создает новый enrich-aware логгер с дополнительными полями
func With(fields ...zap.Field) *WrapLogger {
	if globalLogger == nil {
		return &WrapLogger{zapLogger: zap.NewNop()}
	}

	return &WrapLogger{
		zapLogger: globalLogger.zapLogger.With(fields...),
	}
}

// WithContext создает enrich-aware логгер с контекстом
func WithContext(ctx context.Context) *WrapLogger {
	if globalLogger == nil {
		return &WrapLogger{zapLogger: zap.NewNop()}
	}

	return &WrapLogger{
		zapLogger: globalLogger.zapLogger.With(fieldsFromContext(ctx)...),
	}
}

// Debug - enrich-aware log method
func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.Debug(ctx, msg, fields...)
}

// Info enrich-aware info log
func Info(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.Info(ctx, msg, fields...)
}

// Warn enrich-aware warn log
func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.Warn(ctx, msg, fields...)
}

// Error enrich-aware error log
func Error(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.Error(ctx, msg, fields...)
}

// Fatal enrich-aware fatal log
func Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.Fatal(ctx, msg, fields...)
}

// Методы для обогащения логгера новыми полями из контекста

func (l *WrapLogger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Debug(msg, allFields...)
}

func (l *WrapLogger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Info(msg, allFields...)
}

func (l *WrapLogger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Warn(msg, allFields...)
}

func (l *WrapLogger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Error(msg, allFields...)
}

func (l *WrapLogger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Fatal(msg, allFields...)
}

// Sync сбрасывает буферы логгера
func Sync() error {
	if globalLogger != nil {
		return globalLogger.zapLogger.Sync()
	}
	return nil
}

// buildProductionEncoderConfig - создает достаточный конфиг для необходимого формата логов
func buildProductionEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "timestamp",                 // время
		LevelKey:       "level",                     // уровень логирования
		NameKey:        "logger",                    // имя логгера, если используется
		CallerKey:      "caller",                    // откуда вызван лог
		MessageKey:     "message",                   // текст сообщения
		StacktraceKey:  "stacktrace",                // стектрейс для ошибок
		LineEnding:     zapcore.DefaultLineEnding,   // перенос строки
		EncodeLevel:    zapcore.CapitalLevelEncoder, // INFO, ERROR
		EncodeTime:     zapcore.ISO8601TimeEncoder,  // читаемый ISO 8601 формат
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // короткий caller
		EncodeName:     zapcore.FullNameEncoder,
	}
}

// parseLevel - string to zapcore.Level
func parseLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

func fieldsFromContext(ctx context.Context) []zap.Field {
	fields := make([]zap.Field, 0)

	if traceID, ok := ctx.Value(traceIDKey).(string); ok && traceID != "" {
		fields = append(fields, zap.String(string(traceIDKey), traceID))
	}

	if userID, ok := ctx.Value(userIDKey).(string); ok && userID != "" {
		fields = append(fields, zap.String(string(userIDKey), userID))
	}

	return fields
}
