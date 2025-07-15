package pins

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPanic(t *testing.T) {
	defer func() {
		r := recover()
		assert.NotNil(t, r)
	}()
	PanicOnError(nil)
	PanicOnError(errors.New(""))
}

func TestAssertValue(t *testing.T) {
	defer func() {
		r := recover()
		assert.NotNil(t, r)
	}()
	AssertValue("1", "1")
	AssertValue("1", "2")
}
