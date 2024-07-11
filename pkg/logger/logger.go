package logger

import (
	"context"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	l *logrus.Logger
}

var log Logger

func SetupLogger(l *logrus.Logger) {
	log.l = l
}

func GetLogger() Logger {
	return log
}

func Debug(args ...interface{}) {
	log.l.Debug(args...)
}

func Info(args ...interface{}) {
	log.l.Info(args...)
}

func Warn(args ...interface{}) {
	log.l.Warn(args...)
}

func Error(args ...interface{}) {
	log.l.Error(args...)
}

func Fatal(args ...interface{}) {
	log.l.Fatal(args...)
}

func DebugFn(fn logrus.LogFunction) {
	log.l.DebugFn(fn)
}

func InfoFn(fn logrus.LogFunction) {
	log.l.InfoFn(fn)
}

func WarnFn(fn logrus.LogFunction) {
	log.l.WarnFn(fn)
}

func ErrorFn(fn logrus.LogFunction) {
	log.l.ErrorFn(fn)
}

func FatalFn(fn logrus.LogFunction) {
	log.l.FatalFn(fn)
}

func Debugln(args ...interface{}) {
	log.l.Debug(args...)
}

func Infoln(args ...interface{}) {
	log.l.Infoln(args...)
}

func Warnln(args ...interface{}) {
	log.l.Warnln(args...)
}

func Errorln(args ...interface{}) {
	log.l.Error(args...)
}

func Fatalln(args ...interface{}) {
	log.l.Fatal(args...)
}

func WithField(key string, value interface{}) *logrus.Entry {
	return log.l.WithField(key, value)
}

func WithFields(fields logrus.Fields) *logrus.Entry {
	return log.l.WithFields(fields)
}

func WithError(err error) *logrus.Entry {
	return log.l.WithError(err)
}

func WithContext(ctx context.Context) *logrus.Entry {
	return log.l.WithContext(ctx)
}
