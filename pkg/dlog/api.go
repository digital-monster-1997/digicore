package dlog

// DigitCore TODO 改成吃指令的
var DigitCore = Config{
	Level: "info",
	AddCaller: false,
	Prefix: "DigitCore",
	Debug: true,
	CallerSkip: 0,
}.Build()


func Info(msg string, fields ...Field) {
	DigitCore.Info(msg, fields...)
}

func Debug(msg string, fields ...Field) {
	DigitCore.Debug(msg, fields...)
}

func Warn(msg string, fields ...Field) {
	DigitCore.Warn(msg, fields...)
}

func Error(msg string, fields ...Field) {
	DigitCore.Error(msg, fields...)
}

func Panic(msg string, fields ...Field) {
	DigitCore.Panic(msg, fields...)
}

func DPanic(msg string, fields ...Field) {
	DigitCore.DPanic(msg, fields...)
}

func Fatal(msg string, fields ...Field) {
	DigitCore.Fatal(msg, fields...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	DigitCore.Debugw(msg, keysAndValues...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	DigitCore.Infow(msg, keysAndValues...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	DigitCore.Warnw(msg, keysAndValues...)
}

// Errorw ...
func Errorw(msg string, keysAndValues ...interface{}) {
	DigitCore.Errorw(msg, keysAndValues...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	DigitCore.Panicw(msg, keysAndValues...)
}

func DPanicw(msg string, keysAndValues ...interface{}) {
	DigitCore.DPanicw(msg, keysAndValues...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	DigitCore.Fatalw(msg, keysAndValues...)
}

func Debugf(msg string, args ...interface{}) {
	DigitCore.Debugf(msg, args...)
}

func Infof(msg string, args ...interface{}) {
	DigitCore.Infof(msg, args...)
}

func Warnf(msg string, args ...interface{}) {
	DigitCore.Warnf(msg, args...)
}

func Errorf(msg string, args ...interface{}) {
	DigitCore.Errorf(msg, args...)
}

func Panicf(msg string, args ...interface{}) {
	DigitCore.Panicf(msg, args...)
}

func DPanicf(msg string, args ...interface{}) {
	DigitCore.DPanicf(msg, args...)
}

func Fatalf(msg string, args ...interface{}) {
	DigitCore.Fatalf(msg, args...)
}

func With(fields ...Field) *Logger {
	return DigitCore.With(fields...)
}
