package pins

import (
	"bytes"
	"fmt"

	"github.com/wizards-0/go-pins/logger"
)

func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func LogOnError(err error) {
	if err != nil {
		logger.Error(err)
	}
}

func AssertValue(expected any, actual any) {
	if expected != actual {
		panic(fmt.Errorf("value mismatch, expected: %s, actual: %s", expected, actual))
	}
}

func AppendIfPresent(base *bytes.Buffer, val any, s string) {
	if sVal, ok := val.(string); ok && sVal == "" {
		return
	}
	if iVal, ok := val.(int); ok && iVal == 0 {
		return
	}
	base.Write([]byte(s))
}
