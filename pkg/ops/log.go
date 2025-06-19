package ops

import (
	"fmt"
	"log"
	"log/slog"
	"sync"

	"cloud.google.com/go/logging"
)

type Severity logging.Severity

const (
	Debug    = Severity(logging.Debug)
	Info     = Severity(logging.Info)
	Warning  = Severity(logging.Warning)
	Error    = Severity(logging.Error)
	Critical = Severity(logging.Critical)
)

func (s Severity) Level() slog.Level {
	switch s {
	case Debug:
		return slog.LevelDebug
	case Info:
		return slog.LevelInfo
	case Warning:
		return slog.LevelWarn
	case Error:
		return slog.LevelError
	case Critical:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func (s Severity) String() string {
	switch s {
	case Debug:
		return "DEBUG"
	case Info:
		return "INFO"
	case Warning:
		return "WARNING"
	case Error:
		return "ERROR"
	case Critical:
		return "CRITICAL"
	default:
		return "INFO"
	}
}

var logBuffer = make(chan logging.Entry, 100)

func logFlusher() {
	for entry := range logBuffer {
		log.Printf("[%s] %+v", entry.Severity, entry.Payload)
	}
}

var startLogFlusherOnce sync.Once

func StartLogFlusher() {
	startLogFlusherOnce.Do(func() {
		go logFlusher()
	})
}

// Log logs payload, payload can be string or can marshal into json
func Log(severity Severity, payload interface{}) {
	e := logging.Entry{
		Severity: logging.Severity(severity),
		Payload:  payload,
	}

	if logWriter != nil {
		logWriter.Log(e)
	}

	if logToStd {
		select {
		case logBuffer <- e:
		default:
		}
	}
}

// Logf logs string
func Logf(severity Severity, s string, v ...interface{}) {
	Log(severity, fmt.Sprintf(s, v...))
}
