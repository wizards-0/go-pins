package semver

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wizards-0/go-pins/logger"
)

func TestCompareSemver(t *testing.T) {
	assert := assert.New(t)
	logger.SetLogLevel(logger.LOG_LEVEL_DEBUG)
	verArr := []string{
		"1",
		"2",
		"1.2",
		"1.1",
		"1.11",
		"1.0",
		"2",
		"1.2-alpha",
	}
	sort.Slice(verArr, func(i1, i2 int) bool {
		return CompareSemver(verArr[i1], verArr[i2])
	})

	sortedVerArr := []string{
		"1",
		"1.0",
		"1.1",
		"1.2",
		"1.11",
		"1.2-alpha",
		"2",
		"2",
	}

	assert.Equal(sortedVerArr, verArr)
}
