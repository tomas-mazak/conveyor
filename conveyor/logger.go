package conveyor

import (
	"fmt"
	"log"
)

type Logger struct {
	Ch chan string
}

func (l *Logger) Log(msg string) {
	l.Ch <- msg
}

func (l *Logger) LogFatal(err error) {
	log.Fatalf("CONVEYOR FATAL: %v\n", err)
}

func (l *Logger) LogError(err error) {
	l.Ch <- fmt.Sprintf("CONVEYOR ERROR: %v\n", err)
}
