package http

import (
	"database/sql"
	"github.com/Cococtel/Cococtel_BaBackend/internal/controllers"
	"github.com/gin-gonic/gin"
)

type Router interface {
	MapRoutes()
}
type router struct {
	eng *gin.Engine
	rg  *gin.RouterGroup
	db  *sql.DB
}

func (r *router) MapRoutes() {
	r.setGroup()
	r.buildRoutes()
}
func (r *router) setGroup() {
	r.rg = r.eng.Group("/v1")
}
func (r *router) buildRoutes() {
	r.addSystemPaths()
	//All routes
}
func (r *router) addSystemPaths() {
	r.rg.GET("/ping", controllers.Ping())
}
