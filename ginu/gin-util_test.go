package ginu

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	assert.Equal(http.StatusOK, w.Code)
	assert.Equal("v1", w.Body.String())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", path, strings.NewReader("[\"version\":\"v1\"]"))
	srv.ServeHTTP(w, req)
	assert.Equal(http.StatusBadRequest, w.Code)
}

func TestResponseHandler(t *testing.T) {
	assert := assert.New(t)

	// Success with JSON response
	srv := gin.Default()
	srv.POST(path, BodyHandler(newMigration, func(c *gin.Context, reqBody types.Migration) {
		SendResponseWithStatus(c, http.StatusAccepted, reqBody, nil)
	}))
	m := types.Migration{Version: "v1"}
	json, _ := json.Marshal(m)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", path, bytes.NewReader(json))
	srv.ServeHTTP(w, req)
	assert.Equal(http.StatusAccepted, w.Code)
	assert.Equal(json, w.Body.Bytes())

	// Success with String response
	srv = gin.Default()
	srv.POST(path, BodyHandler(newMigration, func(c *gin.Context, reqBody types.Migration) {
		SendResponse(c, "done", nil)
	}))
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", path, bytes.NewReader(json))
	srv.ServeHTTP(w, req)
	assert.Equal(http.StatusOK, w.Code)
	assert.Equal("done", w.Body.String())

	// Auth error
	srv = gin.Default()
	srv.POST(path, BodyHandler(newMigration, func(c *gin.Context, reqBody types.Migration) {
		SendResponse(c, "done", fmt.Errorf("Unauthorized user"))
	}))
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", path, bytes.NewReader(json))
	srv.ServeHTTP(w, req)
	assert.Equal(http.StatusUnauthorized, w.Code)
	assert.Equal("Unauthorized user", w.Body.String())

	// Not Found
	srv = gin.Default()
	srv.POST(path, BodyHandler(newMigration, func(c *gin.Context, reqBody types.Migration) {
		SendResponse(c, "done", fmt.Errorf("entity not found"))
	}))
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", path, bytes.NewReader(json))
	srv.ServeHTTP(w, req)
	assert.Equal(http.StatusNotFound, w.Code)
	assert.Equal("entity not found", w.Body.String())

	// Random error
	srv = gin.Default()
	srv.POST(path, BodyHandler(newMigration, func(c *gin.Context, reqBody types.Migration) {
		SendResponse(c, "done", fmt.Errorf("server cooked"))
	}))
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", path, bytes.NewReader(json))
	srv.ServeHTTP(w, req)
	assert.Equal(http.StatusInternalServerError, w.Code)
	assert.Equal("server cooked", w.Body.String())
}

func insertMigration(c *gin.Context, reqBody types.Migration) {
	SendResponse(c, reqBody.Version, nil)
}

func newMigration() types.Migration {
	return types.Migration{}
}
