package main

import (
	"log"
	"time"
	"net/http"
)

const (
	ERROR = "ERROR"
	INFO  = "INFO"
	DEBUG = "DEBUG"
)

var LogLevelPriority = map[string]int{
	ERROR: 2,
	INFO:  1,
	DEBUG: 0,
}

type Logger struct {
	Level string // ERROR, INFO, DEBUG
}

func (l Logger) shouldLogLevel(level string) bool {
	return LogLevelPriority[l.Level] <= LogLevelPriority[level]
}

func (l Logger) LogRequest(r *http.Request) {
	if l.shouldLogLevel(INFO) {
		log.Printf(
			"[%s] [INFO] [%s] %s",
			time.Now().UTC(),
			r.Method,
			r.RequestURI,
		)
	}
}

func (l Logger) Error(endpoint string, m string) {
	if l.shouldLogLevel(ERROR) {
		log.Printf(
			"[%s] [ERROR] [%s] %s",
			time.Now(),
			endpoint,
			m,
		)
	}
}

func (l Logger) Info(endpoint string, m string) {
	if l.shouldLogLevel(INFO) {
		log.Printf(
			"[%s] [INFO] [%s] %s",
			time.Now(),
			endpoint,
			m,
		)
	}
}

func (l Logger) Debug(endpoint string, m string) {
	if l.shouldLogLevel(DEBUG) {
		log.Printf(
			"[%s] [DEBUG] [%s] %s",
			time.Now(),
			endpoint,
			m,
		)
	}
}
