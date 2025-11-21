package ginu

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wizards-0/go-pins/logger"
)

func BodyHandler[T any](bodyFactory func() T, fn func(c *gin.Context, reqBody T)) func(c *gin.Context) {
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

func SendResponse(c *gin.Context, resp any, err error) {
	SendResponseWithStatus(c, http.StatusOK, resp, err)
}

func SendResponseWithStatus(c *gin.Context, status int, resp any, err error) {
	if err != nil {
		if strings.Contains(err.Error(), "auth") {
			status = http.StatusUnauthorized
		} else {
			status = http.StatusInternalServerError
		}
		c.String(status, err.Error())
	} else {
		if rs, ok := resp.(string); ok {
			c.String(status, rs)
		} else {
			c.JSON(status, resp)
		}

	}
}
