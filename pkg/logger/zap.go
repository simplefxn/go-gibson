package logger

import (
	"os"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var zLogger zap.Logger
var atom zap.AtomicLevel

// RegisterLog ...
func RegisterLog() error {
	var err error
	zLogger, err = initLog()
	if err != nil {
		return errors.Wrap(err, "RegisterLog")
	}
	defer zLogger.Sync()
	zSugarlog := zLogger.Sugar()
	zSugarlog.Info()

	SetLogger(zSugarlog)
	return nil
}

// GetLogger get the configured logger
func GetLogger() *zap.Logger {
	return &zLogger
}

// initLog create logger
func initLog() (zap.Logger, error) {
	config := zap.NewProductionEncoderConfig()

	config.EncodeTime = zapcore.TimeEncoder(func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.UTC().Format("2006-01-02T15:04:05Z0700"))
		// 2019-08-13T04:39:11Z
	})
	config.LineEnding = ""
	encoder := zapcore.NewConsoleEncoder(config)
	atom = zap.NewAtomicLevel()
	logr := zap.New(zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), atom))

	return *logr, nil
}
