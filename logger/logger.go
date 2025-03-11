package logger

import (
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

func stub(value any) {
	//NOOP if log level not applicable
}
func initError(value any) {
	errLog.Println("ERROR: Logger not initialized. Call one of logger's init func before use")
}
func consoleTrace(value any) {
	outLog.Printf("TRACE: %v\n", value)
}
func consoleDebug(value any) {
	outLog.Printf("DEBUG: %v\n", value)
}
func consoleInfo(value any) {
	outLog.Printf("INFO: %v\n", value)
}
func consoleError(value any) {
	errLog.Printf("ERROR: %v\n", value)
}

var Trace func(any) = initError
var Debug func(any) = initError
var Info func(any) = initError
var Error func(any) = initError

var hasInitialized = false
var level LogLevel = LOG_LEVEL_NONE
var outLog *log.Logger = log.New(os.Stdout, "", log.LstdFlags)
var errLog *log.Logger = log.New(os.Stderr, "", log.LstdFlags)

func Init(logLevel LogLevel) {
	if !hasInitialized {
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

		hasInitialized = true
	} else {
		Error("Logger already initialized with level " + GetLogLevel())
	}
}

func GetLogLevel() string {
	return levelNames[level]
}
