package dlog

import (
	"fmt"
	"github.com/digital-monster-1997/digicore/pkg/utils/dcolor"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"runtime"
)

const (
	// DebugLevel 在 production 環境中禁用，因為數量非常龐大
	DebugLevel = zap.DebugLevel
	// InfoLevel 預設的等級
	InfoLevel = zap.InfoLevel
	// WarnLevel 比 Info 更重要，但不需要團隊開會討論才能設定這個 level，下面等級請通過討論再定義
	WarnLevel = zap.WarnLevel
	// ErrorLevel 較高的優先處理順序，如果應用程式沒啥大錯，不應該跑這個等級
	ErrorLevel = zap.ErrorLevel
	// PanicLevel 記錄一個訊息後 panic
	PanicLevel = zap.PanicLevel
	// FatalLevel 記錄一條訊息，然後應用程式呼叫 os.Exit(1).
	FatalLevel = zap.FatalLevel
)

type Level = zapcore.Level

type Field  = zap.Field

var (
	String = zap.String
	Any = zap.Any
	Int64 = zap.Int64
	Int = zap.Int
	Int32 = zap.Int32
	Uint = zap.Uint
	Duration = zap.Duration
	Durationp = zap.Durationp
	Object = zap.Object
	Namespace = zap.Namespace
	Reflect = zap.Reflect
	Skip = zap.Skip()
	ByteString = zap.ByteString
)

type Logger struct {
	config 	Config
	sugar 	*zap.SugaredLogger
	lv 		*zap.AtomicLevel
	desugar *zap.Logger
}

// IsDebugMode ...
func (logger *Logger) IsDebugMode() bool {
	return logger.config.Debug
}

func normalizeMessage(msg string) string {
	return fmt.Sprintf("%-32s", msg)
}

func sprintf(template string, args ...interface{}) string {
	msg := template
	if msg == "" && len(args) > 0 {
		msg = fmt.Sprint(args...)
	} else if msg != "" && len(args) > 0 {
		msg = fmt.Sprintf(template, args...)
	}
	return msg
}

// StdLog ...
func (logger *Logger) StdLog() *log.Logger {
	return zap.NewStdLog(logger.desugar)
}


//-------------------------- 一般的 --------------------------

func (logger *Logger) Debug(msg string, fields ...Field) {
	if logger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	logger.desugar.Debug(msg, fields...)
}

func (logger *Logger) Info(msg string, fields ...Field) {
	if logger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	logger.desugar.Info(msg, fields...)
}

func (logger *Logger) Warn(msg string, fields ...Field) {
	if logger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	logger.desugar.Warn(msg, fields...)
}

func (logger *Logger) Error(msg string, fields ...Field) {
	if logger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	logger.desugar.Error(msg, fields...)
}

func (logger *Logger) Panic(msg string, fields ...Field) {
	if logger.IsDebugMode() {
		panicDetail(msg, fields...)
		msg = normalizeMessage(msg)
	}
	logger.desugar.Panic(msg, fields...)
}

func (logger *Logger) DPanic(msg string, fields ...Field) {
	if logger.IsDebugMode() {
		panicDetail(msg, fields...)
		msg = normalizeMessage(msg)
	}
	logger.desugar.DPanic(msg, fields...)
}

func (logger *Logger) Fatal(msg string, fields ...Field) {
	if logger.IsDebugMode() {
		panicDetail(msg, fields...)
		msg = normalizeMessage(msg)
		return
	}
	logger.desugar.Fatal(msg, fields...)
}

//-------------------------- 格式化的 --------------------------

func (logger *Logger) Debugf(template string, args ...interface{}) {
	logger.sugar.Debugw(sprintf(template, args...))
}

func (logger *Logger) Infof(template string, args ...interface{}) {
	logger.sugar.Infof(sprintf(template, args...))
}

func (logger *Logger) Warnf(template string, args ...interface{}) {
	logger.sugar.Warnf(sprintf(template, args...))
}

func (logger *Logger) Errorf(template string, args ...interface{}) {
	logger.sugar.Errorf(sprintf(template, args...))
}

func (logger *Logger) Panicf(template string, args ...interface{}) {
	logger.sugar.Panicf(sprintf(template, args...))
}

func (logger *Logger) DPanicf(template string, args ...interface{}) {
	logger.sugar.DPanicf(sprintf(template, args...))
}

func (logger *Logger) Fatalf(template string, args ...interface{}) {
	logger.sugar.Fatalf(sprintf(template, args...))
}
//-------------------------- 記錄一些附加上下文的消息。可變參數鍵值對在 With 中的處理方式相同 --------------------------

func (logger *Logger) Debugw(msg string, keysAndValues ...interface{}) {
	if logger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	logger.sugar.Debugw(msg, keysAndValues...)
}

func (logger *Logger) Infow(msg string, keysAndValues ...interface{}) {
	if logger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	logger.sugar.Infow(msg, keysAndValues...)
}

func (logger *Logger) Warnw(msg string, keysAndValues ...interface{}) {
	if logger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	logger.sugar.Warnw(msg, keysAndValues...)
}

func (logger *Logger) Errorw(msg string, keysAndValues ...interface{}) {
	if logger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	logger.sugar.Errorw(msg, keysAndValues...)
}

func (logger *Logger) Panicw(msg string, keysAndValues ...interface{}) {
	if logger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	logger.sugar.Panicw(msg, keysAndValues...)
}

func (logger *Logger) DPanicw(msg string, keysAndValues ...interface{}) {
	if logger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	logger.sugar.DPanicw(msg, keysAndValues...)
}

func (logger *Logger) Fatalw(msg string, keysAndValues ...interface{}) {
	if logger.IsDebugMode() {
		msg = normalizeMessage(msg)
	}
	logger.sugar.Fatalw(msg, keysAndValues...)
}
func panicDetail(msg string, fields ...Field) {
	enc := zapcore.NewMapObjectEncoder()
	for _, field := range fields {
		field.AddTo(enc)
	}

	// 控制台输出
	fmt.Printf("%s: \n    %s: %s\n", dcolor.Red("panic"), dcolor.Red("msg"), msg)
	if _, file, line, ok := runtime.Caller(3); ok {
		fmt.Printf("    %s: %s:%d\n", dcolor.Red("loc"), file, line)
	}
	for key, val := range enc.Fields {
		fmt.Printf("    %s: %s\n", dcolor.Red(key), fmt.Sprintf("%+v", val))
	}

}

// With ...
func (logger *Logger) With(fields ...Field) *Logger {
	desugarLogger := logger.desugar.With(fields...)
	return &Logger{
		desugar: desugarLogger,
		lv:      logger.lv,
		sugar:   desugarLogger.Sugar(),
		config:  logger.config,
	}
}
