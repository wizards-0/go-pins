package pins

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wizards-0/go-pins/logger"
)

func TestPanic(t *testing.T) {
	defer func() {
		r := recover()
		assert.NotNil(t, r)
	}()
	PanicOnError(nil)
	PanicOnError(errors.New(""))
}

func TestLog(t *testing.T) {
	assert := assert.New(t)
	buf := bytes.Buffer{}
	logger.SetWriter(&buf, &buf, &buf, &buf)
	LogOnError(nil)
	assert.Equal("", buf.String())
	LogOnError(errors.New("test"))
	assert.Contains(buf.String(), "test")
	logger.ResetWriters()
}

func TestAssertValue(t *testing.T) {
	defer func() {
		r := recover()
		assert.NotNil(t, r)
	}()
	AssertValue("1", "1")
	AssertValue("1", "2")
}

func TestAppend(t *testing.T) {
	assert := assert.New(t)
	buf := &bytes.Buffer{}

	someStr := "someString"

	buf.Reset()
	sVal := "test value"
	AppendIfPresent(buf, sVal, someStr)
	assert.Equal("someString", buf.String())

	buf.Reset()
	sVal = ""
	AppendIfPresent(buf, sVal, someStr)
	assert.Equal("", buf.String())

	buf.Reset()
	iVal := 5
	AppendIfPresent(buf, iVal, someStr)
	assert.Equal("someString", buf.String())

	buf.Reset()
	iVal = 0
	AppendIfPresent(buf, iVal, someStr)
	assert.Equal("", buf.String())

	buf.Reset()
	bVal := false
	AppendIfPresent(buf, bVal, someStr)
	assert.Equal("someString", buf.String())

}

func TestMergeError(t *testing.T) {
	assert := assert.New(t)
	assert.Nil(MergeErrors(nil))
	assert.Nil(MergeErrors())
	e1 := fmt.Errorf("msg1")
	e2 := fmt.Errorf("msg2")
	e3 := fmt.Errorf("msg3")
	assert.ErrorContains(MergeErrors(e1), "msg1")
	assert.ErrorContains(MergeErrors(e1, e2), "msg1. msg2")
	assert.ErrorContains(MergeErrors(e1, e2, e3), "msg1. msg2. msg3")
	assert.ErrorContains(MergeErrors(e1, nil), "msg1")
	assert.ErrorContains(MergeErrors(nil, e1, nil), "msg1")
	assert.ErrorContains(MergeErrors(e1, nil, e3), "msg1. msg3")
}
