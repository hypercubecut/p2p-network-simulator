package internal

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func logInitOnce(d bool, f *os.File) {
	onceLogger.Do(func() {
		pe := zap.NewProductionEncoderConfig()

		fileEncoder := zapcore.NewJSONEncoder(pe)

		pe.EncodeTime = zapcore.ISO8601TimeEncoder
		consoleEncoder := zapcore.NewConsoleEncoder(pe)

		level := zap.InfoLevel
		if d {
			level = zap.DebugLevel
		}

		core := zapcore.NewTee(
			zapcore.NewCore(fileEncoder, zapcore.AddSync(f), level),
			zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
		)

		logger = zap.New(core)
	})
}
