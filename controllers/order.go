package controllers

import (
	"backend-technoscape/initializers"
	"backend-technoscape/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateOrder(c *gin.Context) {
	userId, _ := c.Get("userId")
	timeNow := time.Now()
	order1 := models.Order{
		UserID:    userId.(string),
		Amount:    100,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}
	order2 := models.Order{
		UserID:    userId.(string),
		Amount:    200,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}
	initializers.DB.Create(&order1)
	initializers.DB.Create(&order2)
	var user models.User
	initializers.DB.Preload("Orders").Where("id = ? ", userId).First(&user)

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Success create order",
		"data":    user,
	})
	return
}
