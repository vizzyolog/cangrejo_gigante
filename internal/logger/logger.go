package logger

import (
	"fmt"
	"runtime"

	"github.com/sirupsen/logrus"
)

const skipFrames = 2

type LogrusAdapter struct {
	entry *logrus.Entry
}

func New() *LogrusAdapter {
	baseLogger := logrus.New()
	baseLogger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat:   "",
		DisableTimestamp:  false,
		DisableHTMLEscape: false,
		DataKey:           "",
		FieldMap:          nil,
		CallerPrettyfier:  nil,
		PrettyPrint:       true,
	})
	baseLogger.SetLevel(logrus.DebugLevel)

	return &LogrusAdapter{
		entry: logrus.NewEntry(baseLogger),
	}
}

func (l *LogrusAdapter) Info(args ...interface{}) {
	l.entry.WithField("caller", getCaller(skipFrames)).Info(args...)
}

func (l *LogrusAdapter) Infof(format string, args ...interface{}) {
	l.entry.WithField("caller", getCaller(skipFrames)).Infof(format, args...)
}

func (l *LogrusAdapter) Warn(args ...interface{}) {
	l.entry.Warn(args...)
}

func (l *LogrusAdapter) Warnf(format string, args ...interface{}) {
	l.entry.Warnf(format, args...)
}

func (l *LogrusAdapter) Error(args ...interface{}) {
	l.entry.WithField("caller", getCaller(skipFrames)).Error(args...)
}

func (l *LogrusAdapter) Errorf(format string, args ...interface{}) {
	l.entry.WithField("caller", getCaller(skipFrames)).Errorf(format, args...)
}

func (l *LogrusAdapter) WithField(key string, value interface{}) *LogrusAdapter {
	return &LogrusAdapter{
		entry: l.entry.WithField(key, value),
	}
}

func (l *LogrusAdapter) WithFields(fields map[string]interface{}) *LogrusAdapter {
	return &LogrusAdapter{
		entry: l.entry.WithFields(fields),
	}
}

// nolint
func getCaller(skip int) string {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown"
	}

	funcName := runtime.FuncForPC(pc).Name()

	return fmt.Sprintf("%s:%s:%d", file, funcName, line)
}
