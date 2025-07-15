package pins

import "fmt"

func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func AssertValue(expected any, actual any) {
	if expected != actual {
		panic(fmt.Errorf("value mismatch, expected: %s, actual: %s", expected, actual))
	}
}
