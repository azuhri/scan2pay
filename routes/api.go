package routes

import (
	"backend-technoscape/controllers"
	"backend-technoscape/middleware"
	sub "backend-technoscape/routes/sub"

	"github.com/gin-gonic/gin"
)

func ApiRoute(route *gin.Engine) {
	v1 := route.Group("/api/v1")

	auth := v1.Group("/auth")
	{
		auth.POST("/login", controllers.Login)
		auth.POST("/signup", controllers.SignUp)
	}

	app := v1.Group("/app")
	app.Use(middleware.Auth)
	{
		// USER ROUTES
		// userRoutes := app.Group("/user")
		// {
		// 	userRoutes.GET("/", controllers.GetUser)
		// 	userRoutes.PUT("/", controllers.UpdateUser)
		// }
		sub.UserRoutes(app.Group("/user"))
		sub.TransactionRoutes(app.Group("/scan-payment"))
	}
}
