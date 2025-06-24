package logger

import (
	"fmt"
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

const TRACE_PREFIX string = "TRACE: "
const DEBUG_PREFIX string = "DEBUG: "
const INFO_PREFIX string = "INFO: "
const ERROR_PREFIX string = "ERROR: "

var levelNames = [5]string{"Trace", "Debug", "Info", "Error", "None"}

func stub(value ...any) {
	//NOOP if log level not applicable
}

func logErrorStub(err error) error {
	return err
}

func wrapErrorStub(err error, msg string) error {
	return err
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
	b := []byte{}
	b = fmt.Append(b, values...)
	errLog.Output(2, string(b))
}
func consoleLogError(err error) error {
	errLog.Output(2, err.Error())
	return err
}
func checkAndLogErrorToConsole(err error) error {
	if err != nil {
		errLog.Output(2, err.Error())
	}
	return err
}
func wrapAndLogErrorToConsole(err error, msg string) error {
	wrappedError := fmt.Errorf(msg+"\n%w", err)
	errLog.Output(2, wrappedError.Error())
	return wrappedError
}

var Trace func(...any) = stub
var Debug func(...any) = stub
var Info func(...any) = consoleInfo
var Error func(...any) = consoleError
var LogError func(error) error = consoleLogError
var CheckAndLogError func(error) error = checkAndLogErrorToConsole
var WrapAndLogError func(error, string) error = wrapAndLogErrorToConsole

var level LogLevel = LOG_LEVEL_INFO
var traceLog *log.Logger = log.New(os.Stdout, TRACE_PREFIX, log.LstdFlags|log.Lmsgprefix)
var debugLog *log.Logger = log.New(os.Stdout, DEBUG_PREFIX, log.LstdFlags|log.Lmsgprefix)
var infoLog *log.Logger = log.New(os.Stdout, INFO_PREFIX, log.LstdFlags|log.Lmsgprefix)
var errLog *log.Logger = log.New(os.Stderr, ERROR_PREFIX, log.LstdFlags|log.Lmsgprefix|log.Lshortfile)

func SetWriter(
	traceWriter io.Writer,
	debugWriter io.Writer,
	infoWriter io.Writer,
	errWriter io.Writer,
) {
	traceLog = log.New(traceWriter, TRACE_PREFIX, log.LstdFlags|log.Lmsgprefix)
	debugLog = log.New(debugWriter, DEBUG_PREFIX, log.LstdFlags|log.Lmsgprefix)
	infoLog = log.New(infoWriter, INFO_PREFIX, log.LstdFlags|log.Lmsgprefix)
	errLog = log.New(errWriter, ERROR_PREFIX, log.LstdFlags|log.Lmsgprefix|log.Lshortfile)
}

func ResetWriters() {
	traceLog = log.New(os.Stdout, TRACE_PREFIX, log.LstdFlags|log.Lmsgprefix)
	debugLog = log.New(os.Stdout, DEBUG_PREFIX, log.LstdFlags|log.Lmsgprefix)
	infoLog = log.New(os.Stdout, INFO_PREFIX, log.LstdFlags|log.Lmsgprefix)
	errLog = log.New(os.Stderr, ERROR_PREFIX, log.LstdFlags|log.Lmsgprefix|log.Lshortfile)
}

func SetLogLevel(logLevel LogLevel) {

	level = logLevel

	if logLevel <= LOG_LEVEL_ERROR {
		Error = consoleError
		LogError = consoleLogError
		CheckAndLogError = checkAndLogErrorToConsole
		WrapAndLogError = wrapAndLogErrorToConsole
	} else {
		Error = stub
		LogError = logErrorStub
		CheckAndLogError = logErrorStub
		WrapAndLogError = wrapErrorStub
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
