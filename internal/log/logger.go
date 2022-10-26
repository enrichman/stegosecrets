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

func NewSimpleLogger(debug bool) *SimpleLogger {
	return &SimpleLogger{debugEnabled: debug}
}

func (s *SimpleLogger) Print(a ...any) {
	fmt.Println(a...)
}

func (s *SimpleLogger) Debug(a ...any) {
	if s.debugEnabled {
		fmt.Println(a...)
	}
}
