package dlog

import (
	"github.com/digital-monster-1997/digicore/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type Config struct {
	Level 		string
	Fields 		[]zap.Field
	AddCaller 	bool
	Prefix 		string
	Core        zapcore.Core
	Debug       bool
	CallerSkip   int
	EncoderConfig *zapcore.EncoderConfig
}

// RawConfig  讀取 Config 當中的資料
func RawConfig(key string, config *Config, cfg *config.Configuration) *Config {
	if err := cfg.ReadToStruct(key, &config); err != nil {
		panic(err)
	}
	return config
}

// Build ...
func (config Config) Build() *Logger {
	if config.EncoderConfig == nil {
		config.EncoderConfig = DefaultZapConfig()
	}
	if config.Debug {
		config.EncoderConfig.EncodeLevel = DebugEncodeLevel
	}
	logger := newLogger(&config)
	return logger
}

// DefaultZapConfig ...
func DefaultZapConfig() *zapcore.EncoderConfig {
	return &zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "lv",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func newLogger(config *Config) *Logger {
	zapOptions := make([]zap.Option, 0)
	zapOptions = append(zapOptions, zap.AddStacktrace(zap.DPanicLevel))
	if config.AddCaller {
		zapOptions = append(zapOptions, zap.AddCaller(), zap.AddCallerSkip(config.CallerSkip))
	}
	if len(config.Fields) > 0 {
		zapOptions = append(zapOptions, zap.Fields(config.Fields...))
	}
	var ws zapcore.WriteSyncer
	ws = os.Stdout

	lv := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	if err := lv.UnmarshalText([]byte(config.Level)); err != nil {
		panic(err)
	}

	encoderConfig := *config.EncoderConfig
	core := config.Core
	if core == nil {
		core = zapcore.NewCore(
			func() zapcore.Encoder {
				if config.Debug {
					return zapcore.NewConsoleEncoder(encoderConfig)
				}
				return zapcore.NewJSONEncoder(encoderConfig)
			}(),
			ws,
			lv,
		)
	}

	zapLogger := zap.New(
		core,
		zapOptions...,
	)
	return &Logger{
		desugar: zapLogger,
		lv:      &lv,
		config:  *config,
		sugar:   zapLogger.Sugar(),
	}
}
