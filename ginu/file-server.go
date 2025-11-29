package ginu

import (
	"crypto/sha256"
	"encoding/hex"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/wizards-0/go-pins/logger"
)

type staticFile struct {
	data     []byte
	mimeType string
	eTag     string
}

var cacheMutex sync.RWMutex
var fileCache = map[string]staticFile{}

func File(ctx *gin.Context, filePath string, cacheControl string) {
	cacheMutex.RLock()
	file, found := fileCache[filePath]
	cacheMutex.RUnlock()
	if !found {
		cacheMutex.Lock()
		fileData, err := os.ReadFile(filePath)
		if err != nil {
			err = logger.WrapAndLogError(err, "error in reading file to be served, path: "+filePath)
			ctx.String(http.StatusNotFound, err.Error())
			cacheMutex.Unlock()
			return
		}
		ext := filepath.Ext(filePath)
		mimeType := mime.TypeByExtension(ext)
		sha256Hash := sha256.Sum256(fileData)
		file = staticFile{
			data:     fileData,
			mimeType: mimeType,
			eTag:     hex.EncodeToString(sha256Hash[:]),
		}
		fileCache[filePath] = file
		cacheMutex.Unlock()
	}
	if cacheControl == "" {
		cacheControl = "no-cache"
	}
	ctx.Header("Cache-Control", cacheControl)
	ctx.Header("Etag", file.eTag)
	if ctx.GetHeader("If-None-Match") != file.eTag {
		ctx.Header("Content-Type", file.mimeType)
		ctx.Writer.Write(file.data)
	} else {
		ctx.Status(http.StatusNotModified)
	}
}

func ClearFileServerCache() {
	clear(fileCache)
}
