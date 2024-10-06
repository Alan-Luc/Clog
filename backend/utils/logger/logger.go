package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitLogger() {
	var err error
	Logger, err = zap.NewDevelopment(zap.AddStacktrace(zapcore.ErrorLevel))
	if err != nil {
		panic(err)
	}
}
