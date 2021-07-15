package logger

import (
	"fmt"
)

type Logger struct {
}

func (l Logger) Write(bytes []byte) (int, error) {
	return fmt.Print(string(bytes))
}
