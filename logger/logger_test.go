package logger

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var buf = bytes.Buffer{}
var logReader = io.Reader(&buf)
var testMsg = "Test log"

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
func getErrorLog(msg string) string {
	return getDateString() + " ERROR: " + msg + "\n"
}

func TestLoggerInit(t *testing.T) {
	assert := assert.New(t)
	setup()
	assert.Equal("Info", GetLogLevel())
	Trace(testMsg)
	logMsg, _ := io.ReadAll(logReader)
	assert.Equal("", string(logMsg))
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
	assert.Equal(getErrorLog(testMsg), string(logMsg))
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
	assert.Equal(getErrorLog(testMsg), string(logMsg))
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
	assert.Equal(getErrorLog(testMsg), string(logMsg))
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
	assert.Equal(getErrorLog(testMsg), string(logMsg))
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
}
