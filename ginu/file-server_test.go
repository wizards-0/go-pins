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
	// Valid File, default cache-control
	srv := gin.Default()
	srv.GET(path, func(ctx *gin.Context) {
		File(ctx, filePath, "")
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", path, nil)
	srv.ServeHTTP(w, req)
	respCc := w.Result().Header["Cache-Control"]
	assert.Equal(http.StatusOK, w.Code)
	assert.Equal("no-cache", respCc[0])
	etag := w.Result().Header["Etag"]
	assert.NotEmpty(etag)

	// Valid File with if-none-match header
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", path, nil)
	req.Header["If-None-Match"] = etag
	srv.ServeHTTP(w, req)
	assert.Equal(http.StatusNotModified, w.Code)

	// Valid File with if-none-match header, after cache reset
	ClearFileServerCache()
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", path, nil)
	req.Header["If-None-Match"] = []string{"f1f133f"}
	srv.ServeHTTP(w, req)
	assert.Equal(http.StatusOK, w.Code)

	// Valid File Second time, custom cache-control
	w = httptest.NewRecorder()
	srv = gin.Default()
	srv.GET(path, func(ctx *gin.Context) {
		File(ctx, filePath, "max-age=60")
	})
	req, _ = http.NewRequest("GET", path, nil)
	srv.ServeHTTP(w, req)
	respCc = w.Result().Header["Cache-Control"]
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
