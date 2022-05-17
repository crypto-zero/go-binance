package common

import (
	"log"
	"testing"
)

func TestLogPrint(t *testing.T) {
	logger := NewDefaultLogger(LogDebug, log.Default())
	logger.Debugw("hello world")
	logger.Debugw("hello world", "a", "b")
	logger.Debugw("hello world", "a")
	logger.Debugw("hello world", "a", "b", "c", map[string]string{"hello": "x"})
	logger.Debugw("hello world", "a", "b", "d", struct{ Test string }{Test: "dev"})
	logger.Debugw("hello world", "a", "b", "d", &struct{ Test string }{Test: "prod"})
}
