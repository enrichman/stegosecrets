package log

import (
	"fmt"
	"io"
)

type Level int

const (
	None Level = iota
	Info
	Debug
)

func NewLevel(silent, verbose bool) Level {
	if silent {
		return None
	} else if verbose {
		return Debug
	}

	return Info
}

type Logger interface {
	Print(a ...any)
	Debug(a ...any)
}

type SimpleLogger struct {
	writer io.Writer
	level  Level
}

func NewSimpleLogger(writer io.Writer, level Level) *SimpleLogger {
	if level == None {
		writer = io.Discard
	}

	return &SimpleLogger{writer: writer, level: level}
}

func (s *SimpleLogger) Print(a ...any) {
	if s.level >= Info {
		fmt.Fprintln(s.writer, a...)
	}
}

func (s *SimpleLogger) Debug(a ...any) {
	if s.level >= Debug {
		fmt.Fprintln(s.writer, a...)
	}
}
