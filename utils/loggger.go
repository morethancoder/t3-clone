package utils

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type Logger interface {
	Info(format string, args ...any)
	Success(format string, args ...any)
	Error(format string, args ...any)
	Loading(format string, args ...any)
	Warn(format string, args ...any)
	Debug(format string, args ...any)
	Fatal(format string, args ...any)
}

type logMode string

const (
	infoMode    logMode = "INFO"
	successMode logMode = "SUCCESS"
	errorMode   logMode = "ERROR"
	loadingMode logMode = "LOADING"
	warnMode    logMode = "WARN"
	debugMode   logMode = "DEBUG"
	fatalMode   logMode = "FATAL"
)

const (
	Reset       = "\033[0m"
	Primary     = "\033[33m" // Yellow (used for warnings)
	Secondary   = "\033[36m" // Cyan (used for debug)
	Accent      = "\033[37m" // White (used for loading)
	InfoColor   = "\033[96m" // Bright Cyan (used for info)
	Success     = "\033[32m" // Green (used for success)
	Warning     = "\033[93m" // Bright Yellow (used for warnings)
	Error       = "\033[31m" // Red (used for errors)
	FatalColor  = "\033[35m" // Magenta (used for fatal)
	Dim         = "\033[2m"  // Dim text
)

type logger struct {
	Out       io.Writer
	DebugMode bool
	mutex     sync.Mutex
}

var Log *logger

func init() {
	Log = NewLogger()
}

func NewLogger() *logger {
	return &logger{
		Out:       os.Stdout,
	}
}

func (l *logger) prefix(mode logMode, color string) string {
	timeNow := fmt.Sprintf("%s%s%s", Accent, time.Now().Format(time.Kitchen), Reset)
	return fmt.Sprintf("[%s] %s%s%s ", timeNow, color, mode, Reset)
}

func (l *logger) log(mode logMode, color string, format string, args ...any) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	message := fmt.Sprintf(format, args...)
	_, err := l.Out.Write([]byte(l.prefix(mode, color) + Dim + message + Reset + "\n"))
	if err != nil {
		_, _ = fmt.Printf("Logging error: %v\n", err)
	}
}

func (l *logger) Info(format string, args ...any) {
	l.log(infoMode, InfoColor, format, args...)
}

func (l *logger) Success(format string, args ...any) {
	l.log(successMode, Success, format, args...)
}

func (l *logger) Error(format string, args ...any) {
	l.log(errorMode, Error, format, args...)
}

func (l *logger) Loading(format string, args ...any) {
	l.log(loadingMode, Accent, format, args...)
}

func (l *logger) Warn(format string, args ...any) {
	l.log(warnMode, Warning, format, args...)
}

func (l *logger) Debug(format string, args ...any) {
	if l.DebugMode {
		l.log(debugMode, Secondary, format, args...)
	}
}

func (l *logger) Fatal(format string, args ...any) {
	l.log(fatalMode, FatalColor, format, args...)
	os.Exit(1)
}
