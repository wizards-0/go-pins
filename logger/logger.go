package logger

import (
	"io"
	"log"
	"os"
)

type LogLevel int

const (
	LOG_LEVEL_TRACE LogLevel = iota
	LOG_LEVEL_DEBUG LogLevel = iota
	LOG_LEVEL_INFO  LogLevel = iota
	LOG_LEVEL_ERROR LogLevel = iota
	LOG_LEVEL_NONE  LogLevel = iota
)

var levelNames = [5]string{"Trace", "Debug", "Info", "Error", "None"}

func stub(value ...any) {
	//NOOP if log level not applicable
}

func consoleTrace(values ...any) {
	traceLog.Println(values...)
}
func consoleDebug(values ...any) {
	debugLog.Println(values...)
}
func consoleInfo(values ...any) {
	infoLog.Println(values...)
}
func consoleError(values ...any) {
	errLog.Println(values...)
}

var Trace func(...any) = stub
var Debug func(...any) = stub
var Info func(...any) = consoleInfo
var Error func(...any) = consoleError

var level LogLevel = LOG_LEVEL_INFO
var traceLog *log.Logger = log.New(os.Stdout, "TRACE: ", log.LstdFlags|log.Lmsgprefix)
var debugLog *log.Logger = log.New(os.Stdout, "DEBUG: ", log.LstdFlags|log.Lmsgprefix)
var infoLog *log.Logger = log.New(os.Stdout, "INFO: ", log.LstdFlags|log.Lmsgprefix)
var errLog *log.Logger = log.New(os.Stderr, "ERROR: ", log.LstdFlags|log.Lmsgprefix)

func SetWriter(
	traceWriter io.Writer,
	debugWriter io.Writer,
	infoWriter io.Writer,
	errWriter io.Writer,
) {
	traceLog = log.New(traceWriter, "TRACE: ", log.LstdFlags|log.Lmsgprefix)
	debugLog = log.New(debugWriter, "DEBUG: ", log.LstdFlags|log.Lmsgprefix)
	infoLog = log.New(infoWriter, "INFO: ", log.LstdFlags|log.Lmsgprefix)
	errLog = log.New(errWriter, "ERROR: ", log.LstdFlags|log.Lmsgprefix)
}

func SetLogLevel(logLevel LogLevel) {

	level = logLevel

	if logLevel <= LOG_LEVEL_ERROR {
		Error = consoleError
	} else {
		Error = stub
	}

	if logLevel <= LOG_LEVEL_INFO {
		Info = consoleInfo
	} else {
		Info = stub
	}

	if logLevel <= LOG_LEVEL_DEBUG {
		Debug = consoleDebug
	} else {
		Debug = stub
	}

	if logLevel <= LOG_LEVEL_TRACE {
		Trace = consoleTrace
	} else {
		Trace = stub
	}
}

func GetLogLevel() string {
	return levelNames[level]
}
