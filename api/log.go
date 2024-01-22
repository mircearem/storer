package api

import (
	"io"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	_ "github.com/mircearem/storer/log"
	"github.com/sirupsen/logrus"
)

type apiLogger struct {
	log *logrus.Logger
}

func NewApiLogger(log *logrus.Logger) echo.Logger {
	return &apiLogger{log: log}
}

// Output implements the echo.Logger interface
func (l *apiLogger) Output() io.Writer {
	return l.log.Out
}

// SetOutput implements the echo.Logger interface
func (l *apiLogger) SetOutput(w io.Writer) {
	l.log.SetOutput(w)
}

// Prefix implements the echo.Logger interface
func (l *apiLogger) Prefix() string {
	return ""
}

// SetPrefix implements the echo.Logger interface
func (l *apiLogger) SetPrefix(p string) {}

// Level implements the echo.Logger interface
func (l *apiLogger) Level() log.Lvl {
	switch l.log.GetLevel() {
	case logrus.DebugLevel:
		return log.DEBUG
	case logrus.InfoLevel:
		return log.INFO
	case logrus.WarnLevel:
		return log.WARN
	case logrus.ErrorLevel:
		return log.ERROR
	case logrus.FatalLevel:
		return log.ERROR
	case logrus.PanicLevel:
		return log.ERROR
	}
	return log.INFO
}

// SetLevel implements the echo.Logger interface
func (l *apiLogger) SetLevel(level log.Lvl) {
	l.log.SetLevel(logrus.Level(level))
}

// SetHeader is not used by logrus
func (l *apiLogger) SetHeader(s string) {}

// Print implements the echo.Logger interface
func (l *apiLogger) Print(i ...interface{}) {
	l.log.Print(i...)
}

// Printf implements the echo.Logger interface
func (l *apiLogger) Printf(format string, args ...interface{}) {
	l.log.Printf(format, args...)
}

// Printj implements the echo.Logger interface
func (l *apiLogger) Printj(j log.JSON) {
	l.log.WithFields(logrus.Fields(j)).Print()
}

// Print implements the echo.Logger interface
func (l *apiLogger) Debug(i ...interface{}) {
	l.log.Debug(i...)
}

// Debugf implements the echo.Logger interface
func (l *apiLogger) Debugf(format string, args ...interface{}) {
	l.log.Debugf(format, args...)
}

// Debugj implements the echo.Logger interface
func (l *apiLogger) Debugj(j log.JSON) {
	l.log.WithFields(logrus.Fields(j)).Debug()
}

// Info implements the echo.Logger interface
func (l *apiLogger) Info(i ...interface{}) {
	l.log.Info(i...)
}

// Infof implements the echo.Logger interface
func (l *apiLogger) Infof(format string, args ...interface{}) {
	l.log.Infof(format, args...)
}

// Infoj implements the echo.Logger interface
func (l *apiLogger) Infoj(j log.JSON) {
	l.log.WithFields(logrus.Fields(j)).Info()
}

// Warn implements the echo.Logger interface
func (l *apiLogger) Warn(i ...interface{}) {
	l.log.Warn(i...)
}

// Warn implements the echo.Logger interface
func (l *apiLogger) Warnf(format string, args ...interface{}) {
	l.log.Warnf(format, args...)
}

// Warnj implements the echo.Logger interface
func (l *apiLogger) Warnj(j log.JSON) {
	l.log.WithFields(logrus.Fields(j)).Warn()
}

// Error implements the echo.Logger interface
func (l *apiLogger) Error(i ...interface{}) {
	l.log.Error(i...)
}

// Errorf implements the echo.Logger interface
func (l *apiLogger) Errorf(format string, args ...interface{}) {
	l.log.Errorf(format, args...)
}

// Errorj implements the echo.Logger interface
func (l *apiLogger) Errorj(j log.JSON) {
	l.log.WithFields(logrus.Fields(j)).Error()
}

// Fatal implements the echo.Logger interface
func (l *apiLogger) Fatal(i ...interface{}) {
	l.log.Fatal(i...)
}

// Fatalf implements the echo.Logger interface
func (l *apiLogger) Fatalf(format string, args ...interface{}) {
	l.log.Fatalf(format, args...)
}

// Fatalj implements the echo.Logger interface
func (l *apiLogger) Fatalj(j log.JSON) {
	l.log.WithFields(logrus.Fields(j)).Fatal()
}

// Panic implements the echo.Logger interface
func (l *apiLogger) Panic(i ...interface{}) {
	l.log.Panic(i...)
}

// Panicf implements the echo.Logger interface
func (l *apiLogger) Panicf(format string, args ...interface{}) {
	l.log.Panicf(format, args...)
}

// Panicj implements the echo.Logger interface
func (l *apiLogger) Panicj(j log.JSON) {
	l.log.WithFields(logrus.Fields(j)).Panic()
}
