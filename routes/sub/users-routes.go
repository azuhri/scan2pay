package routes

import (
	"backend-technoscape/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.RouterGroup) {
	r.GET("/", controllers.GetUser)
	r.PUT("/", controllers.UpdateUser)
	r.POST("/bank-account/activation", controllers.ActivationBankAccount)

	routePin := r.Group("/pin")
	{
		routePin.POST("", controllers.SetupPin)
		routePin.POST("/verify", controllers.VerifyPinNumber)
	}
}
