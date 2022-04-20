package dlog

import (
	"github.com/digital-monster-1997/digicore/pkg/utils/dcolor"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// DebugEncodeLevel ...
func DebugEncodeLevel(lv zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	var colorize = dcolor.Red
	switch lv {
	case zapcore.DebugLevel:
		colorize = dcolor.Blue
	case zapcore.InfoLevel:
		colorize = dcolor.Green
	case zapcore.WarnLevel:
		colorize = dcolor.Yellow
	case zapcore.ErrorLevel, zap.PanicLevel, zap.DPanicLevel, zap.FatalLevel:
		colorize = dcolor.Red
	default:
	}
	enc.AppendString(colorize(lv.CapitalString()))
}
