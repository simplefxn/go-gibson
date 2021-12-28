package logger

import "go.uber.org/zap"

// SaramaLogger wraps zap.Logger with the sarama.StdLogger interface
type SaramaLogger struct {
	logger *zap.SugaredLogger
}

// NewSaramaLogger initializes a new SaramaLogger with a zap.Logger
func NewSaramaLogger(zl *zap.Logger) *SaramaLogger {
	return &SaramaLogger{
		logger: zl.Sugar(),
	}
}

func (s *SaramaLogger) Print(v ...interface{}) {
	s.logger.Info(v...)
}

func (s *SaramaLogger) Printf(format string, v ...interface{}) {
	s.logger.Infof(format, v...)
}

func (s *SaramaLogger) Println(v ...interface{}) {
	s.logger.Info(v...)
}
