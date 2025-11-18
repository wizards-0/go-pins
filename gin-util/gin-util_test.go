package ginu

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/wizards-0/go-pins/migrator/types"
)

var path = "/testReqBody"

func TestBodyHandler(t *testing.T) {
	assert := assert.New(t)
	srv := gin.Default()
	srv.POST(path, BodyHandler(newMigration, insertMigration))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", path, strings.NewReader("{\"version\":\"v1\"}"))
	srv.ServeHTTP(w, req)
	assert.Equal(http.StatusAccepted, w.Code)
	assert.Equal("v1", w.Body.String())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", path, strings.NewReader("[\"version\":\"v1\"]"))
	srv.ServeHTTP(w, req)
	assert.Equal(http.StatusBadRequest, w.Code)
}

func insertMigration(c *gin.Context, reqBody types.Migration) {
	c.String(http.StatusAccepted, reqBody.Version)
}

func newMigration() types.Migration {
	return types.Migration{}
}
