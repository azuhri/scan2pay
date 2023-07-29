package routes

import (
	"backend-technoscape/controllers"

	"github.com/gin-gonic/gin"
)

func OrderRoutes(r *gin.RouterGroup) {
	r.POST("/", controllers.CreateOrder)
}
