package http

import (
	"database/sql"
	"github.com/gin-gonic/gin"
)

func InitRouter(eng *gin.Engine, db *sql.DB) Router {
	return &router{eng: eng, db: db}
}
