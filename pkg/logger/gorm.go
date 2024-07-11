package logger

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
)

// LogMode sets the log mode for the GormLogger
func (l Logger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

// Info logs an informational message
func (l Logger) Info(ctx context.Context, s string, i ...interface{}) {
	log.l.WithContext(ctx).Infof(s, i...)
}

// Warn logs a warning message
func (l Logger) Warn(ctx context.Context, s string, i ...interface{}) {
	log.l.WithContext(ctx).Warnf(s, i...)
}

// Error logs an error message
func (l Logger) Error(ctx context.Context, s string, i ...interface{}) {
	log.l.WithContext(ctx).Errorf(s, i...)
}

// Trace logs a SQL statement with its duration
func (l Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if err != nil {
		sql, rows := fc()
		log.l.WithContext(ctx).WithError(err).WithFields(logrus.Fields{
			"sql":        sql,
			"rows":       rows,
			"elapsed_ms": time.Since(begin).Milliseconds(),
		}).Error("SQL execution failed")
	} else {
		sql, rows := fc()
		log.l.WithContext(ctx).WithFields(logrus.Fields{
			"sql":        sql,
			"rows":       rows,
			"elapsed_ms": time.Since(begin).Milliseconds(),
		}).Info("SQL executed")
	}
}
