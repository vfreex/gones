package logger

import (
	"go.uber.org/zap"
)

type Logger interface {
	Debug(msg string)
	Info(msg string)
}

var logger *ZapLogger

func init() {
	logger = NewLogger()
}

func GetLogger() *ZapLogger {
	return logger
}

type ZapLogger struct {
	logger *zap.Logger
	suger  *zap.SugaredLogger
}

func NewLogger() *ZapLogger {
	logger, _ := zap.NewProduction()
	//logger, _ = zap.NewDevelopment()
	return &ZapLogger{
		logger: logger,
		suger:  logger.Sugar(),
	}
}

func (l *ZapLogger) Debug(msg string) {
	//l.logger.Debug(msg)
}

func (l *ZapLogger) Debugf(msg string, args ...interface{}) {
	//l.suger.Debugf(msg, args...)
}

func (l *ZapLogger) Info(msg string) {
	l.logger.Info(msg)
}

func (l *ZapLogger) Infof(msg string, args ...interface{}) {
	l.suger.Infof(msg, args...)
}

func (l *ZapLogger) Warn(msg string) {
	l.logger.Warn(msg)
}

func (l *ZapLogger) Warnf(msg string, args ...interface{}) {
	l.suger.Warnf(msg, args...)
}

func (l *ZapLogger) Fatal(msg string) {
	l.logger.Fatal(msg)
}

func (l *ZapLogger) Fatalf(msg string, args ...interface{}) {
	l.suger.Fatalf(msg, args...)
}

func (l *ZapLogger) Sync() {
	l.logger.Sync()
}