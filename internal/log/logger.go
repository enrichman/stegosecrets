package log

import (
	"fmt"
)

type Logger interface {
	Print(a ...any)
	Debug(a ...any)
}

type SimpleLogger struct {
	debugEnabled bool
}

type SilentLogger struct {
}

func NewSimpleLogger(debug bool) *SimpleLogger {
	return &SimpleLogger{debugEnabled: debug}
}

func NewSilentLogger() *SilentLogger {
	return &SilentLogger{}
}

func (s *SimpleLogger) Print(a ...any) {
	fmt.Println(a...)
}

func (s *SimpleLogger) Debug(a ...any) {
	if s.debugEnabled {
		fmt.Println(a...)
	}
}

func (s *SilentLogger) Print(a ...any) {
	// SilentLogger - Do nothing
}

func (s *SilentLogger) Debug(a ...any) {
	// SilentLogger - Do nothing
}
