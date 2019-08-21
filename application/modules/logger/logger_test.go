package logger

import (
	"bytes"
	"fmt"
	"github.com/onsi/gomega"
	"strings"
	"testing"
)

const (
	TesteLogDebug     = "Teste de print de debug"
	TesteLogInfo      = "Teste de print de info"
	TesteLogNotice    = "Teste de print de notice"
	TesteLogWarning   = "Teste de print de warning"
	TesteLogError     = "Teste de print de error"
	TesteLogCritical  = "Teste de print de critical"
	TesteLogAlert     = "Teste de print de alert"
	TesteLogEmergency = "Teste de print de emergency"
)

func configLog(Out *bytes.Buffer, level int) {
	Testlog := GetLogger()
	Testlog.SetLogOutput(Out)
	Testlog.SetLogFlags(0)
	Testlog.SetLogLevel(level)
}

func logPrintMessages() {
	Testlog := GetLogger()
	Testlog.DebugPrintln(TesteLogDebug)
	Testlog.InfoPrintln(TesteLogInfo)
	Testlog.NoticePrintln(TesteLogNotice)
	Testlog.WarningPrintln(TesteLogWarning)
	Testlog.ErrorPrintln(TesteLogError)
	Testlog.CriticalPrintln(TesteLogCritical)
	Testlog.AlertPrintln(TesteLogAlert)
	Testlog.EmergencyPrintln(TesteLogEmergency)
}

func TestLogPrint(test *testing.T) {
	gomega.RegisterTestingT(test)
	Out := bytes.NewBuffer(nil)
	configLog(Out, LogDebug)
	logPrintMessages()
	gomega.Expect(strings.Contains(fmt.Sprint(Out), TesteLogDebug)).To(gomega.BeTrue())
}

func TestLogDebugLevelPrint(test *testing.T) {
	gomega.RegisterTestingT(test)

	Out := bytes.NewBuffer(nil)
	configLog(Out, LogDebug)
	logPrintMessages()
	OutMockString := ""
	OutMockString += fmt.Sprintln(TesteLogDebug)
	OutMockString += fmt.Sprintln(TesteLogInfo)
	OutMockString += fmt.Sprintln(TesteLogNotice)
	OutMockString += fmt.Sprintln(TesteLogWarning)
	OutMockString += fmt.Sprintln(TesteLogError)
	OutMockString += fmt.Sprintln(TesteLogCritical)
	OutMockString += fmt.Sprintln(TesteLogAlert)
	OutMockString += fmt.Sprintln(TesteLogEmergency)

	gomega.Expect(strings.Compare(fmt.Sprint(Out), OutMockString)).To(gomega.Equal(0))
}

func TestLogInfoLevelPrint(test *testing.T) {
	gomega.RegisterTestingT(test)

	Out := bytes.NewBuffer(nil)
	configLog(Out, LogInfo)
	logPrintMessages()
	OutMockString := ""
	OutMockString += fmt.Sprintln(TesteLogInfo)
	OutMockString += fmt.Sprintln(TesteLogNotice)
	OutMockString += fmt.Sprintln(TesteLogWarning)
	OutMockString += fmt.Sprintln(TesteLogError)
	OutMockString += fmt.Sprintln(TesteLogCritical)
	OutMockString += fmt.Sprintln(TesteLogAlert)
	OutMockString += fmt.Sprintln(TesteLogEmergency)

	gomega.Expect(strings.Compare(fmt.Sprint(Out), OutMockString)).To(gomega.Equal(0))
}

func TestLogNoticeLevelPrint(test *testing.T) {
	gomega.RegisterTestingT(test)

	Out := bytes.NewBuffer(nil)
	configLog(Out, LogNotice)
	logPrintMessages()
	OutMockString := ""
	OutMockString += fmt.Sprintln(TesteLogNotice)
	OutMockString += fmt.Sprintln(TesteLogWarning)
	OutMockString += fmt.Sprintln(TesteLogError)
	OutMockString += fmt.Sprintln(TesteLogCritical)
	OutMockString += fmt.Sprintln(TesteLogAlert)
	OutMockString += fmt.Sprintln(TesteLogEmergency)

	gomega.Expect(strings.Compare(fmt.Sprint(Out), OutMockString)).To(gomega.Equal(0))
}

func TestLogWarningLevelPrint(test *testing.T) {
	gomega.RegisterTestingT(test)

	Out := bytes.NewBuffer(nil)
	configLog(Out, LogWarning)
	logPrintMessages()
	OutMockString := ""
	OutMockString += fmt.Sprintln(TesteLogWarning)
	OutMockString += fmt.Sprintln(TesteLogError)
	OutMockString += fmt.Sprintln(TesteLogCritical)
	OutMockString += fmt.Sprintln(TesteLogAlert)
	OutMockString += fmt.Sprintln(TesteLogEmergency)

	gomega.Expect(strings.Compare(fmt.Sprint(Out), OutMockString)).To(gomega.Equal(0))
}

func TestLogErrorLevelPrint(test *testing.T) {
	gomega.RegisterTestingT(test)

	Out := bytes.NewBuffer(nil)
	configLog(Out, LogError)
	logPrintMessages()

	OutMockString := ""
	OutMockString += fmt.Sprintln(TesteLogError)
	OutMockString += fmt.Sprintln(TesteLogCritical)
	OutMockString += fmt.Sprintln(TesteLogAlert)
	OutMockString += fmt.Sprintln(TesteLogEmergency)

	gomega.Expect(strings.Compare(fmt.Sprint(Out), OutMockString)).To(gomega.Equal(0))
}

func TestLogCriticalLevelPrint(test *testing.T) {
	gomega.RegisterTestingT(test)

	Out := bytes.NewBuffer(nil)
	configLog(Out, LogCritical)
	logPrintMessages()
	OutMockString := ""
	OutMockString += fmt.Sprintln(TesteLogCritical)
	OutMockString += fmt.Sprintln(TesteLogAlert)
	OutMockString += fmt.Sprintln(TesteLogEmergency)

	gomega.Expect(strings.Compare(fmt.Sprint(Out), OutMockString)).To(gomega.Equal(0))
}

func TestLogAlertLevelPrint(test *testing.T) {
	gomega.RegisterTestingT(test)

	Out := bytes.NewBuffer(nil)
	configLog(Out, LogAlert)
	logPrintMessages()
	OutMockString := ""
	OutMockString += fmt.Sprintln(TesteLogAlert)
	OutMockString += fmt.Sprintln(TesteLogEmergency)

	gomega.Expect(strings.Compare(fmt.Sprint(Out), OutMockString)).To(gomega.Equal(0))
}

func TestLogEmergencyLevelPrint(test *testing.T) {
	gomega.RegisterTestingT(test)

	Out := bytes.NewBuffer(nil)
	configLog(Out, LogEmergency)
	logPrintMessages()
	OutMockString := ""
	OutMockString += fmt.Sprintln(TesteLogEmergency)
	gomega.Expect(strings.Compare(fmt.Sprint(Out), OutMockString)).To(gomega.Equal(0))
}
