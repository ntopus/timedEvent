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

var instance *StdLogger
var once sync.Once

type StdLogger struct {
	logger   *log.Logger
	logLevel int
	logDepht int
}

func GetLogger() *StdLogger {
	once.Do(func() {
		instance = &StdLogger{logger: log.New(os.Stderr, "", log.LstdFlags), logLevel: LogNotice, logDepht: 2}
	})
	return instance
}
func (l *StdLogger) SetLogLevel(LogLevel int) {
	l.logLevel = LogLevel
}
func (l *StdLogger) SetLogOutput(out io.Writer) {
	l.logger.SetOutput(out)
}
func (l *StdLogger) SetLogPrefix(prefix string) {
	l.logger.SetPrefix(prefix)
}

func (l *StdLogger) GetLogFlags() int {
	return l.logger.Flags()
}

func (l *StdLogger) SetLogFlags(flags int) {
	l.logger.SetFlags(flags)
}

func (l *StdLogger) DebugPrintln(v ...interface{}) {
	if l.logLevel <= LogDebug {
		l.logger.Output(l.logDepht, fmt.Sprintln(v...))
	}
}
func (l *StdLogger) InfoPrintln(v ...interface{}) {
	if l.logLevel <= LogInfo {
		l.logger.Output(l.logDepht, fmt.Sprintln(v...))
	}
}
func (l *StdLogger) NoticePrintln(v ...interface{}) {
	if l.logLevel <= LogNotice {
		l.logger.Output(l.logDepht, fmt.Sprintln(v...))
	}
}
func (l *StdLogger) WarningPrintln(v ...interface{}) {
	if l.logLevel <= LogWarning {
		l.logger.Output(l.logDepht, fmt.Sprintln(v...))
	}
}
func (l *StdLogger) ErrorPrintln(v ...interface{}) {
	if l.logLevel <= LogError {
		l.logger.Output(l.logDepht, fmt.Sprintln(v...))
	}
}
func (l *StdLogger) CriticalPrintln(v ...interface{}) {
	if l.logLevel <= LogCritical {
		l.logger.Output(l.logDepht, fmt.Sprintln(v...))
	}
}
func (l *StdLogger) AlertPrintln(v ...interface{}) {
	if l.logLevel <= LogAlert {
		l.logger.Output(l.logDepht, fmt.Sprintln(v...))
	}
}
func (l *StdLogger) EmergencyPrintln(v ...interface{}) {
	if l.logLevel <= LogEmergency {
		l.logger.Output(l.logDepht, fmt.Sprintln(v...))
	}
}
