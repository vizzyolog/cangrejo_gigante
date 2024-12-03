package logger

import (
	"fmt"
	"runtime"

	"github.com/sirupsen/logrus"
)

type logrusAdapter struct {
	entry *logrus.Entry
}

func New() Logger {
	baseLogger := logrus.New()
	baseLogger.SetFormatter(&logrus.JSONFormatter{})
	baseLogger.SetLevel(logrus.DebugLevel)

	return &logrusAdapter{
		entry: logrus.NewEntry(baseLogger),
	}
}

func (l *logrusAdapter) Info(args ...interface{}) {
	l.entry.WithField("caller", getCaller(2)).Info(args...)
}

func (l *logrusAdapter) Infof(format string, args ...interface{}) {
	l.entry.WithField("caller", getCaller(2)).Infof(format, args...)
}

func (l *logrusAdapter) Warn(args ...interface{}) {
	l.entry.WithField("caller", getCaller(2)).Warn(args...)
}

func (l *logrusAdapter) Warnf(format string, args ...interface{}) {
	l.entry.WithField("caller", getCaller(2)).Warnf(format, args...)
}

func (l *logrusAdapter) Error(args ...interface{}) {
	l.entry.WithField("caller", getCaller(2)).Error(args...)
}

func (l *logrusAdapter) Errorf(format string, args ...interface{}) {
	l.entry.WithField("caller", getCaller(2)).Errorf(format, args...)
}

func (l *logrusAdapter) WithField(key string, value interface{}) Logger {
	return &logrusAdapter{
		entry: l.entry.WithField(key, value),
	}
}

func (l *logrusAdapter) WithFields(fields map[string]interface{}) Logger {
	return &logrusAdapter{
		entry: l.entry.WithFields(fields),
	}
}

func getCaller(skip int) string {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown"
	}
	funcName := runtime.FuncForPC(pc).Name()
	return file + ":" + funcName + ":" + fmt.Sprintf("%d", line)
}
