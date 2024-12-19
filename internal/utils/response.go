package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type response struct {
	Data interface{} `json:"data"`
}
type errorResponse struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"string"`
}

func Response(c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}
func Success(c *gin.Context, status int, data interface{}) {
	Response(c, status, response{Data: data})
}
func Error(c *gin.Context, status int, format string, args ...interface{}) {
	err := errorResponse{
		Code:    strings.ReplaceAll(strings.ToLower(http.StatusText(status)), " ", "_"),
		Message: fmt.Sprintf(format, args),
		Status:  status,
	}
	Response(c, status, err)
}
