package ginu

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

/*

	2. Invalid File
	3. Valid File with if-none-match header
	4. Valid File second time
*/

func TestFileServe(t *testing.T) {
	assert := assert.New(t)
	filePath := "../resources/test/properties/common.properties"
	// Valid File
	srv := gin.Default()
	srv.GET(path, func(ctx *gin.Context) {
		File(ctx, filePath, "")
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", path, nil)
	srv.ServeHTTP(w, req)
	assert.Equal(http.StatusOK, w.Code)
	etag := w.Result().Header["Etag"]
	assert.NotEmpty(etag)

	// Valid File with if-none-match header
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", path, nil)
	req.Header["If-None-Match"] = etag
	srv.ServeHTTP(w, req)
	assert.Equal(http.StatusNotModified, w.Code)

	// Valid File Second time
	w = httptest.NewRecorder()
	srv = gin.Default()
	srv.GET(path, func(ctx *gin.Context) {
		File(ctx, filePath, "max-age=60")
	})
	req, _ = http.NewRequest("GET", path, nil)
	srv.ServeHTTP(w, req)
	respCc := w.Result().Header["Cache-Control"]
	assert.Equal(http.StatusOK, w.Code)
	assert.Equal("max-age=60", respCc[0])

	// Invalid file
	srv = gin.Default()
	srv.GET(path, func(ctx *gin.Context) {
		File(ctx, "invalid-path", "")
	})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", path, nil)
	srv.ServeHTTP(w, req)
	assert.Equal(http.StatusNotFound, w.Code)

}
