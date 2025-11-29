package ginu

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wizards-0/go-pins/logger"
)

func BodyHandler[T any](bodyFactory func() T, fn func(ctx *gin.Context, reqBody T)) func(c *gin.Context) {
	body := bodyFactory()
	return func(c *gin.Context) {
		if err := json.NewDecoder(c.Request.Body).Decode(&body); err != nil {
			logger.Error(err)
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		fn(c, body)
	}
}

func SendResponse(ctx *gin.Context, resp any, err error) {
	SendResponseWithStatus(ctx, http.StatusOK, resp, err)
}

func SendResponseWithStatus(ctx *gin.Context, status int, resp any, err error) {
	if err != nil {
		if strings.Contains(err.Error(), "auth") {
			status = http.StatusUnauthorized
		} else if strings.Contains(err.Error(), "not found") {
			status = http.StatusNotFound
		} else {
			status = http.StatusInternalServerError
		}
		ctx.String(status, err.Error())
	} else {
		if rs, ok := resp.(string); ok {
			ctx.String(status, rs)
		} else {
			ctx.JSON(status, resp)
		}

	}
}
