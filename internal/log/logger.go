package logger

import (
	"github.com/sirupsen/logrus"
)

// Logger indicates minimal method to implement logger
type Logger interface {
	Debug(...interface{})
	Info(...interface{})
	Warn(...interface{})
	Error(...interface{})
	Fatal(...interface{})
}

// logger will be an extendable method from logrus.Entry
type logger struct {
	f logrus.Fields
}

func (l *logger) Debug(args ...interface{}) {
	e := logrus.WithFields(l.f)
	e.Debug(args...)
}

func (l *logger) Info(args ...interface{}) {
	e := logrus.WithFields(l.f)
	e.Info(args...)
}

func (l *logger) Warn(args ...interface{}) {
	e := logrus.WithFields(l.f)
	e.Warn(args...)
}

func (l *logger) Error(args ...interface{}) {
	e := logrus.WithFields(l.f)
	e.Error(args...)
}

func (l *logger) Fatal(args ...interface{}) {
	e := logrus.WithFields(l.f)
	e.Fatal(args...)
}
