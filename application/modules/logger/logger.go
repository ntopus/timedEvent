package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

const (
	LogDebug     = 100
	LogInfo      = 200
	LogNotice    = 250
	LogWarning   = 300
	LogError     = 400
	LogCritical  = 500
	LogAlert     = 550
	LogEmergency = 600
)

var instance *stdLogger
var once sync.Once

type stdLogger struct {
	logger   *log.Logger
	logLevel int
	logDepht int
}

func GetLogger() *stdLogger {
	once.Do(func() {
		instance = &stdLogger{logger: log.New(os.Stderr, "", log.LstdFlags), logLevel: LogNotice, logDepht: 2}
	})
	return instance
}
func (l *stdLogger) SetLogLevel(LogLevel int) {
	l.logLevel = LogLevel
}
func (l *stdLogger) SetLogOutput(out io.Writer) {
	l.logger.SetOutput(out)
}
func (l *stdLogger) SetLogPrefix(prefix string) {
	l.logger.SetPrefix(prefix)
}

func (l *stdLogger) GetLogFlags() int {
	return l.logger.Flags()
}

func (l *stdLogger) SetLogFlags(flags int) {
	l.logger.SetFlags(flags)
}

func (l *stdLogger) DebugPrintln(v ...interface{}) {
	if l.logLevel <= LogDebug {
		l.logger.Output(l.logDepht, fmt.Sprintln(v...))
	}
}
func (l *stdLogger) InfoPrintln(v ...interface{}) {
	if l.logLevel <= LogInfo {
		l.logger.Output(l.logDepht, fmt.Sprintln(v...))
	}
}
func (l *stdLogger) NoticePrintln(v ...interface{}) {
	if l.logLevel <= LogNotice {
		l.logger.Output(l.logDepht, fmt.Sprintln(v...))
	}
}
func (l *stdLogger) WarningPrintln(v ...interface{}) {
	if l.logLevel <= LogWarning {
		l.logger.Output(l.logDepht, fmt.Sprintln(v...))
	}
}
func (l *stdLogger) ErrorPrintln(v ...interface{}) {
	if l.logLevel <= LogError {
		l.logger.Output(l.logDepht, fmt.Sprintln(v...))
	}
}
func (l *stdLogger) CriticalPrintln(v ...interface{}) {
	if l.logLevel <= LogCritical {
		l.logger.Output(l.logDepht, fmt.Sprintln(v...))
	}
}
func (l *stdLogger) AlertPrintln(v ...interface{}) {
	if l.logLevel <= LogAlert {
		l.logger.Output(l.logDepht, fmt.Sprintln(v...))
	}
}
func (l *stdLogger) EmergencyPrintln(v ...interface{}) {
	if l.logLevel <= LogEmergency {
		l.logger.Output(l.logDepht, fmt.Sprintln(v...))
	}
}
