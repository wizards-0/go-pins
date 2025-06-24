package migrator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wizards-0/go-pins/logger"
	"github.com/wizards-0/go-pins/migrator/types"
)

const invalid_filename = "invalid filename"

func TestValidDir(t *testing.T) {
	setup()
	assert := assert.New(t)
	logger.SetLogLevel(logger.LOG_LEVEL_DEBUG)
	migrations, err := parseDirectory("../resources/test/migrations/valid")
	assert.Nil(err)
	assert.Equal("user-setup", migrations[0].Name)

	migrations, err = parseDirectory("../resources/test/migrations/valid/")
	assert.Nil(err)
	assert.Equal("user-setup", migrations[0].Name)
}

func TestValidMultiLevelDir(t *testing.T) {
	setup()
	assert := assert.New(t)

	migrations, err := parseDirectory("../resources/test/migrations/valid-multi-level/")
	assert.Nil(err)
	assert.Equal(2, len(migrations))
}

func TestInvalidFilenameFormat(t *testing.T) {
	setup()
	assert := assert.New(t)
	_, err := parseDirectory("../resources/test/migrations/invalid-filename/invalid-format/")
	assert.ErrorContains(err, invalid_filename)
}

func TestInvalidFiletype(t *testing.T) {
	setup()
	assert := assert.New(t)
	_, err := parseDirectory("../resources/test/migrations/invalid-filename/invalid-filetype/")
	assert.ErrorContains(err, invalid_filename)
}

func TestInvalidFilenameQueryType(t *testing.T) {
	setup()
	assert := assert.New(t)
	_, err := parseDirectory("../resources/test/migrations/invalid-filename/invalid-query-type/")
	assert.ErrorContains(err, invalid_filename)
}

func TestNameMismatch(t *testing.T) {
	setup()
	assert := assert.New(t)
	_, err := parseDirectory("../resources/test/migrations/invalid-filename/name-mismatch/")
	assert.ErrorContains(err, invalid_filename)
}

func TestMissingQuery(t *testing.T) {
	setup()
	assert := assert.New(t)
	_, err := parseDirectory("../resources/test/migrations/missing-query/")
	assert.ErrorContains(err, "missing query")
}

func TestMissingRollback(t *testing.T) {
	setup()
	assert := assert.New(t)
	_, err := parseDirectory("../resources/test/migrations/missing-rollback/")
	assert.ErrorContains(err, "missing rollback")
}

func TestInvalidPath(t *testing.T) {
	setup()
	assert := assert.New(t)
	_, err := parseDirectory("../invalid-path/")
	assert.ErrorContains(err, "The system cannot find the file specified")
}

func TestInvalidFilePath(t *testing.T) {
	setup()
	assert := assert.New(t)
	err := addFileToMap("../invalid-path/", "invalid-file.txt", map[string]types.Migration{})
	assert.ErrorContains(err, "The system cannot find the file specified")
}
