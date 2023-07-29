package routes

import (
	"backend-technoscape/controllers"

	"github.com/gin-gonic/gin"
)

func TransactionRoutes(r *gin.RouterGroup) {
	r.POST("/:uuid", controllers.CreateTransaction)
	r.GET("/", controllers.GetListTransaction)
}
