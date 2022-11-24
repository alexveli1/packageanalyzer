// Package mylog is custom logger based on uber/zap
package mylog

import (
	"log"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github/alexveli1/packageanalyzer/internal/domain"
)

var SugarLogger *zap.SugaredLogger

func InitLogger(out string, filename string) *zap.SugaredLogger {
	defer func() {
		if re := recover(); re != nil {
			log.Printf("recovered error in initating logger:%v", re)
		}
	}()
	var writeSyncer zapcore.WriteSyncer
	encoder := getEncoder()
	if strings.EqualFold(out, domain.LogTypeFile) || out == "" {
		f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			log.Printf("error opening log file: %v", err)
		}
		writeSyncer = zapcore.AddSync(f)
	}
	if strings.EqualFold(out, domain.LogTypeStdOut) {
		writeSyncer = zapcore.AddSync(os.Stdout)
	}
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zap.ErrorLevel))
	SugarLogger := logger.Sugar()

	return SugarLogger
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "name",
		CallerKey:      "caller",
		FunctionKey:    "function",
		StacktraceKey:  "stack",
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.RFC3339TimeEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
	}

	return zapcore.NewConsoleEncoder(encoderConfig)
}
