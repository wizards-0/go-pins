package logger

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var buf = bytes.Buffer{}
var logReader = io.Reader(&buf)
var testMsg = "Test log"
var testWrapMsg = "Caller Info"

func setup() {
	w := io.Writer(&buf)
	SetWriter(w, w, w, w)
}

func getDateString() string {
	return time.Now().Format("2006/01/02 15:04:05")
}

func getTraceLog(msg string) string {
	return getDateString() + " TRACE: " + msg + "\n"
}
func getDebugLog(msg string) string {
	return getDateString() + " DEBUG: " + msg + "\n"
}
func getInfoLog(msg string) string {
	return getDateString() + " INFO: " + msg + "\n"
}
func getErrorLog(msg string, lineNumber string) string {
	return getDateString() + " logger_test.go:" + lineNumber + ": ERROR: " + msg + "\n"
}

func TestLoggerInit(t *testing.T) {
	assert := assert.New(t)
	setup()
	assert.Equal("Info", GetLogLevel())
	Trace(testMsg)
	logMsg, _ := io.ReadAll(logReader)
	assert.Equal("", string(logMsg))

	ResetWriters()
	assert.Equal(traceLog.Writer(), os.Stdout)
}

func TestLoggerLevelTrace(t *testing.T) {
	assert := assert.New(t)
	setup()

	SetLogLevel(LOG_LEVEL_TRACE)

	Trace(testMsg)
	logMsg, _ := io.ReadAll(logReader)
	assert.Equal(getTraceLog(testMsg), string(logMsg))

	Debug(testMsg)
	logMsg, _ = io.ReadAll(logReader)
	assert.Equal(getDebugLog(testMsg), string(logMsg))

	Info(testMsg)
	logMsg, _ = io.ReadAll(logReader)
	assert.Equal(getInfoLog(testMsg), string(logMsg))

	Error(testMsg)
	logMsg, _ = io.ReadAll(logReader)
	assert.Equal(getErrorLog(testMsg, "71"), string(logMsg))

	LogError(errors.New(testMsg))
	logMsg, _ = io.ReadAll(logReader)
	assert.Equal(getErrorLog(testMsg, "75"), string(logMsg))

	WrapAndLogError(errors.New(testMsg), testWrapMsg)
	logMsg, _ = io.ReadAll(logReader)
	assert.Equal(getErrorLog(testWrapMsg+"\n"+testMsg, "79"), string(logMsg))
}

func TestLoggerLevelDebug(t *testing.T) {
	assert := assert.New(t)
	setup()

	SetLogLevel(LOG_LEVEL_DEBUG)

	Trace(testMsg)
	logMsg, _ := io.ReadAll(logReader)
	assert.Equal("", string(logMsg))

	Debug(testMsg)
	logMsg, _ = io.ReadAll(logReader)
	assert.Equal(getDebugLog(testMsg), string(logMsg))

	Info(testMsg)
	logMsg, _ = io.ReadAll(logReader)
	assert.Equal(getInfoLog(testMsg), string(logMsg))

	Error(testMsg)
	logMsg, _ = io.ReadAll(logReader)
	assert.Equal(getErrorLog(testMsg, "102"), string(logMsg))

	LogError(errors.New(testMsg))
	logMsg, _ = io.ReadAll(logReader)
	assert.Equal(getErrorLog(testMsg, "106"), string(logMsg))

	WrapAndLogError(errors.New(testMsg), testWrapMsg)
	logMsg, _ = io.ReadAll(logReader)
	assert.Equal(getErrorLog(testWrapMsg+"\n"+testMsg, "110"), string(logMsg))
}

func TestLoggerLevelInfo(t *testing.T) {
	assert := assert.New(t)
	setup()

	SetLogLevel(LOG_LEVEL_INFO)

	Trace(testMsg)
	logMsg, _ := io.ReadAll(logReader)
	assert.Equal("", string(logMsg))

	Debug(testMsg)
	logMsg, _ = io.ReadAll(logReader)
	assert.Equal("", string(logMsg))

	Info(testMsg)
	logMsg, _ = io.ReadAll(logReader)
	assert.Equal(getInfoLog(testMsg), string(logMsg))

	Error(testMsg)
	logMsg, _ = io.ReadAll(logReader)
	assert.Equal(getErrorLog(testMsg, "133"), string(logMsg))

	LogError(errors.New(testMsg))
	logMsg, _ = io.ReadAll(logReader)
	assert.Equal(getErrorLog(testMsg, "137"), string(logMsg))

	WrapAndLogError(errors.New(testMsg), testWrapMsg)
	logMsg, _ = io.ReadAll(logReader)
	assert.Equal(getErrorLog(testWrapMsg+"\n"+testMsg, "141"), string(logMsg))

}

func TestLoggerLevelError(t *testing.T) {
	assert := assert.New(t)
	setup()

	SetLogLevel(LOG_LEVEL_ERROR)

	Trace(testMsg)
	logMsg, _ := io.ReadAll(logReader)
	assert.Equal("", string(logMsg))

	Debug(testMsg)
	logMsg, _ = io.ReadAll(logReader)
	assert.Equal("", string(logMsg))

	Info(testMsg)
	logMsg, _ = io.ReadAll(logReader)
	assert.Equal("", string(logMsg))

	Error(testMsg)
	logMsg, _ = io.ReadAll(logReader)
	assert.Equal(getErrorLog(testMsg, "165"), string(logMsg))

	LogError(errors.New(testMsg))
	logMsg, _ = io.ReadAll(logReader)
	assert.Equal(getErrorLog(testMsg, "169"), string(logMsg))

	WrapAndLogError(errors.New(testMsg), testWrapMsg)
	logMsg, _ = io.ReadAll(logReader)
	assert.Equal(getErrorLog(testWrapMsg+"\n"+testMsg, "173"), string(logMsg))

}

func TestLoggerLevelNone(t *testing.T) {
	assert := assert.New(t)
	setup()

	SetLogLevel(LOG_LEVEL_NONE)

	Trace(testMsg)
	logMsg, _ := io.ReadAll(logReader)
	assert.Equal("", string(logMsg))

	Debug(testMsg)
	logMsg, _ = io.ReadAll(logReader)
	assert.Equal("", string(logMsg))

	Info(testMsg)
	logMsg, _ = io.ReadAll(logReader)
	assert.Equal("", string(logMsg))

	Error(testMsg)
	logMsg, _ = io.ReadAll(logReader)
	assert.Equal("", string(logMsg))

	LogError(errors.New(testMsg))
	logMsg, _ = io.ReadAll(logReader)
	assert.Equal("", string(logMsg))

	WrapAndLogError(errors.New(testMsg), testWrapMsg)
	logMsg, _ = io.ReadAll(logReader)
	assert.Equal("", string(logMsg))
}

func TestErrorLogging(t *testing.T) {
	assert := assert.New(t)
	setup()
	SetLogLevel(LOG_LEVEL_ERROR)

	errMsg := "Source Error"
	e := errors.New(errMsg)
	_ = LogError(e)
	logMsg, _ := io.ReadAll(logReader)
	assert.Equal(getErrorLog(errMsg, "217"), string(logMsg))

	_ = WrapAndLogError(e, testWrapMsg)
	logMsg, _ = io.ReadAll(logReader)
	assert.Equal(getErrorLog(testWrapMsg+"\n"+errMsg, "221"), string(logMsg))

	_ = CheckAndLogError(nil)
	logMsg, _ = io.ReadAll(logReader)
	assert.Equal("", string(logMsg))

	_ = CheckAndLogError(e)
	logMsg, _ = io.ReadAll(logReader)
	assert.Equal(getErrorLog(errMsg, "229"), string(logMsg))
}
