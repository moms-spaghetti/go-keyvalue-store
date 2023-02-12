package logger

import (
	"log"
	"time"
)

type Logger struct {
	logging    chan string
	done       chan struct{}
	timeFormat string
}

func NewLogger() *Logger {
	logging := make(chan string)
	done := make(chan struct{})

	return &Logger{
		logging:    logging,
		done:       done,
		timeFormat: time.RFC822,
	}
}

func (l Logger) Start() {
	log.Print("logger started")
	go func() {
		for {
			select {
			case out := <-l.logging:
				l.print(out)
			case <-l.done:
				return
			}
		}
	}()
}

func (l Logger) StartNoopLogger() {
	go func() {
		for {
			select {
			case <-l.logging:
			case <-l.done:
				return
			}
		}
	}()
}

func (l Logger) Log(s string) {
	l.logging <- s
}

func (l Logger) print(s string) {
	t := time.Now().Format(l.timeFormat)

	log.Printf("Logger @ %s: %s", t, s)
}

func (l Logger) Stop() {
	close(l.logging)
	close(l.done)
	log.Print("logger shutdown ok")
}
