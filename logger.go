package logutil

import (
	"os"
	"sync"
	"time"
	"zaplog/zapencoder"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	gLogger     *zap.Logger
	gLoggerOnce sync.Once
)

// MillisecondDurationEncoder serializes a time.Duration to a floating-point number of seconds elapsed.
func MillisecondDurationEncoder(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendFloat64(float64(d) / float64(time.Millisecond))
}

func Suger() *zap.SugaredLogger {
	return GetLogger().Sugar()
}

func GetLogger() *zap.Logger {
	gLoggerOnce.Do(func() {
		encoderCfg := zapcore.EncoderConfig{
			MessageKey:     "message",
			LevelKey:       "level",
			TimeKey:        "time",
			NameKey:        "logger",
			CallerKey:      "", // "caller",
			StacktraceKey:  "stacktrace",
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.000"),
			EncodeDuration: MillisecondDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}

		encoder := zapencoder.NewTextEncoder(encoderCfg)
		cores := []zapcore.Core{
			zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel),
		}
		core := zapcore.NewTee(cores...)
		gLogger = zap.New(core)

		// config := DefaultZapLoggerConfig
		// lg, err := config.Build()
		// if err != nil {
		// 	panic("s")
		// }
		// gLogger = lg
	})

	return gLogger
}
