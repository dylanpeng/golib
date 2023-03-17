package gorm

import (
	"context"
	oLogger "github.com/dylanpeng/golib/logger"
	"gorm.io/gorm"
	gLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"strings"
	"time"
)

type logger struct {
	logger   *oLogger.Logger
	LogLevel gLogger.LogLevel
}

// LogMode log mode
func (l *logger) LogMode(level gLogger.LogLevel) gLogger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

// Info print info
func (l *logger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gLogger.Info {
		l.logger.Info(msg, data)
	}
}

// Warn print warn messages
func (l *logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gLogger.Warn {
		l.logger.Warn(msg, data)
	}
}

// Error print error messages
func (l *logger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gLogger.Error {
		l.logger.Error(msg, data)
	}
}

// Trace print sql message
func (l *logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= gLogger.Silent {
		return
	}

	useTime := time.Since(begin)
	hasErr := (err != nil && err != gorm.ErrRecordNotFound)

	source := utils.FileWithLineNum()

	if dirs := strings.Split(source, "/"); len(dirs) >= 3 {
		source = strings.Join(dirs[len(dirs)-3:], "/")
	}

	sql, rows := fc()

	if hasErr {
		if rows == -1 {
			l.logger.Infof("query: <%s> | %4v | - | %s | %s", source, useTime, sql, err)
		} else {
			l.logger.Infof("query: <%s> | %4v | %d rows | %s | %s", source, useTime, rows, sql, err)
		}
	} else {
		if rows == -1 {
			l.logger.Infof("query: <%s> | %4v | - | %s", source, useTime, sql)
		} else {
			l.logger.Infof("query: <%s> | %4v | %d rows | %s", source, useTime, rows, sql)
		}
	}

	return
}
