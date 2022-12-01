package routes

import (
	"github.com/gin-gonic/gin"

	"golang-resto-management/controllers"
)

func OrderRoutes(incommingRoutes *gin.Engine) {
	incommingRoutes.GET("/orders", controllers.GetOrders())
	incommingRoutes.GET("/order/:order_id", controllers.GetOrder())
	incommingRoutes.POST("/order", controllers.CreateOrder())
	incommingRoutes.PATCH("/order/:order_id", controllers.EditOrder())
}
