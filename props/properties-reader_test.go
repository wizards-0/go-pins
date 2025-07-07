package props

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadFiles(t *testing.T) {
	assert := assert.New(t)
	props, err := ReadFiles("../resources/test/properties/common.properties")
	assert.Nil(err)
	assert.Equal("8080", props["PORT"])

	props, err = ReadFiles("../resources/test/properties/common.properties", "../resources/test/properties/local.properties")
	assert.Nil(err)
	assert.Equal("4002", props["PORT"])
	assert.Equal("local", props["ENV_NAME"])
	assert.Equal("P22=\\", props["PWD"])
}

func TestReadFilesErrors(t *testing.T) {
	assert := assert.New(t)
	_, err := ReadFiles("../invalid-path")
	assert.ErrorContains(err, "error in reading file")

	_, err = ReadFiles("../resources/test/properties/invalid.properties")
	assert.ErrorContains(err, "invalid property")
}
