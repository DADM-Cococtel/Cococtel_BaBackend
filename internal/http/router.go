package http

import (
	"database/sql"
	"github.com/Cococtel/Cococtel_BaBackend/internal/controllers"
	"github.com/Cococtel/Cococtel_BaBackend/internal/controllers/usercontroller"
	"github.com/Cococtel/Cococtel_BaBackend/internal/defines"
	"github.com/Cococtel/Cococtel_BaBackend/internal/middleware"
	"github.com/Cococtel/Cococtel_BaBackend/internal/repository/userrepository"
	"github.com/Cococtel/Cococtel_BaBackend/internal/services/userservice"
	"github.com/gin-gonic/gin"
)

type Router interface {
	MapRoutes()
}

type router struct {
	eng *gin.Engine
	db  *sql.DB
}

func (r *router) MapRoutes() {
	r.setGroup()
	r.addSystemPaths()
	r.buildRoutes()
}

func (r *router) setGroup() {
	r.eng.Use(middleware.CORS())
	//r.rg = r.eng.Group("/v1", middleware.ProtectedHandler())
}

func (r *router) buildRoutes() {
	userRepository := userrepository.NewUserRepository(r.db)

	userService := userservice.NewUser(userRepository)

	userController := usercontroller.NewUser(userService)

	r.eng.POST(defines.VerifyPath, userController.VerifyUser())
	r.eng.POST(defines.RegisterPath, userController.RegisterUser())
	r.eng.POST(defines.LoginPath, userController.LoginUser())
	r.eng.POST(defines.ValidateLoginPath, userController.ValidateLogin())
	r.eng.POST(defines.GetQRDoubleAuthPath, userController.GetQRDoubleAuth())
	r.eng.POST(defines.NotifyQRReadPath, userController.NotifyQRRead())
}
func (r *router) addSystemPaths() {
	r.eng.GET(defines.PingPath, controllers.Ping())
}
