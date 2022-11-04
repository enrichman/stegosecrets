package log

import (
	"fmt"
	"io"
)

type Logger interface {
	Print(a ...any)
	Debug(a ...any)
}

type SimpleLogger struct {
	writer       io.Writer
	debugEnabled bool
}

func NewSimpleLogger(writer io.Writer, debug bool) *SimpleLogger {
	return &SimpleLogger{writer: writer, debugEnabled: debug}
}

func (s *SimpleLogger) Print(a ...any) {
	fmt.Fprintln(s.writer, a...)
}

func (s *SimpleLogger) Debug(a ...any) {
	if s.debugEnabled {
		fmt.Fprintln(s.writer, a...)
	}
}

// SilentLogger is a nop Logger implementation.
type SilentLogger struct{}

func (s *SilentLogger) Print(a ...any) {}

func (s *SilentLogger) Debug(a ...any) {}
