package coze

import (
	"context"
	"fmt"
	"log"
	"os"
)

// Logger ...
type Logger interface {
	Log(ctx context.Context, level LogLevel, message string, args ...interface{})
}

type LevelLogger interface {
	Logger
	SetLevel(level LogLevel)
}

type LogLevel int

// LogLevelTrace ...
const (
	LogLevelTrace LogLevel = iota + 1
	LogLevelDebug
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

// String ...
func (r LogLevel) String() string {
	switch r {
	case LogLevelTrace:
		return "TRACE"
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	default:
		return ""
	}
}

type stdLogger struct {
	log *log.Logger
}

// newStdLogger ...
func newStdLogger() Logger {
	return &stdLogger{
		log: log.New(os.Stderr, "", log.LstdFlags),
	}
}

// Log ...
func (l *stdLogger) Log(ctx context.Context, level LogLevel, message string, args ...interface{}) {
	if len(args) == 0 {
		_ = l.log.Output(2, "["+level.String()+"] "+message)
	} else {
		_ = l.log.Output(2, "["+level.String()+"] "+fmt.Sprintf(message, args...))
	}
}

type levelLogger struct {
	Logger
	level LogLevel
}

// NewLevelLogger ...
func NewLevelLogger(logger Logger, level LogLevel) LevelLogger {
	return &levelLogger{
		Logger: logger,
		level:  level,
	}
}

// SetLevel ...
func (l *levelLogger) SetLevel(level LogLevel) {
	l.level = level
}

// Log ...
func (l *levelLogger) Log(ctx context.Context, level LogLevel, message string, args ...interface{}) {
	if level >= l.level {
		l.Logger.Log(ctx, level, message, args...)
	}
}

func (l *levelLogger) Debugf(ctx context.Context, message string, args ...interface{}) {
	l.Log(ctx, LogLevelDebug, message, args...)
}

func (l *levelLogger) Infof(ctx context.Context, message string, args ...interface{}) {
	l.Log(ctx, LogLevelInfo, message, args...)
}

func (l *levelLogger) Warnf(ctx context.Context, message string, args ...interface{}) {
	l.Log(ctx, LogLevelWarn, message, args...)
}

func (l *levelLogger) Errorf(ctx context.Context, message string, args ...interface{}) {
	l.Log(ctx, LogLevelError, message, args...)
}

var logger = levelLogger{
	Logger: newStdLogger(),
	level:  LogLevelInfo,
}

func setLogger(l Logger) {
	logger.Logger = l
}

func setLevel(level LogLevel) {
	logger.level = level
}

func (r *core) Log(ctx context.Context, level LogLevel, msg string, args ...interface{}) {
	if level >= r.logLevel {
		if r.logger != nil {
			r.logger.Log(ctx, level, msg, args...)
		} else {
			logger.Log(ctx, level, msg, args...)
		}
	}
}
