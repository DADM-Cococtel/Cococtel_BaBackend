package controllers

import (
	"github.com/Cococtel/Cococtel_BaBackend/internal/utils"
	"github.com/gin-gonic/gin"
)

func Ping() gin.HandlerFunc {
	return func(c *gin.Context) {
		utils.Success(c, 200, "Pong")
	}
}
