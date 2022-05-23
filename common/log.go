package common

import "fmt"

type Logger interface {
	Debugw(msg string, keyAndValues ...interface{})
	Infow(msg string, keyAndValues ...interface{})
	Warningw(msg string, keyAndValues ...interface{})
}

type Printf interface {
	Printf(format string, v ...interface{})
}

type LogLevel int

const (
	LogDebug LogLevel = iota + 1
	LogInfo
	LogWarning
)

type DefaultLogger struct {
	level LogLevel
	p     Printf
}

func NewDefaultLogger(level LogLevel, p Printf) *DefaultLogger {
	return &DefaultLogger{level: level, p: p}
}

func (d *DefaultLogger) SetLevel(level LogLevel) {
	d.level = level
}

func (d *DefaultLogger) levelEnable(level LogLevel) bool {
	return level >= d.level
}

func (d *DefaultLogger) formatConcat(msg string, keyAndValues []interface{}) (format string,
	values []interface{},
) {
	suffixFormat := ""
	for i := 0; i < len(keyAndValues) && i+2 <= len(keyAndValues); i += 2 {
		if i > 0 {
			suffixFormat += ","
		}
		suffixFormat += fmt.Sprintf(" %s: ", keyAndValues[i]) + "%#v"
		values = append(values, keyAndValues[i+1])
	}
	return msg + suffixFormat, values
}

func (d *DefaultLogger) Debugw(msg string, keyAndValues ...interface{}) {
	if !d.levelEnable(LogDebug) {
		return
	}
	msg, values := d.formatConcat(msg, keyAndValues)
	d.p.Printf(msg, values...)
}

func (d *DefaultLogger) Infow(msg string, keyAndValues ...interface{}) {
	if !d.levelEnable(LogInfo) {
		return
	}
	msg, values := d.formatConcat(msg, keyAndValues)
	d.p.Printf(msg, values...)
}

func (d *DefaultLogger) Warningw(msg string, keyAndValues ...interface{}) {
	if !d.levelEnable(LogWarning) {
		return
	}
	msg, values := d.formatConcat(msg, keyAndValues)
	d.p.Printf(msg, values...)
}
