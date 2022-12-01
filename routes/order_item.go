package routes

import (
	"github.com/gin-gonic/gin"

	"golang-resto-management/controllers"
)

func OrderItemRoutes(incommingRoutes *gin.Engine) {
	incommingRoutes.GET("/order_items", controllers.GetOrderItems())
	incommingRoutes.GET("/order_item/:order_item_id", controllers.GetOrderItem())
	incommingRoutes.GET("/order_items/:order_id", controllers.GetOrderItemsByOrder())
	incommingRoutes.POST("/order_item", controllers.CreateOrderItem())
	incommingRoutes.PATCH("/order_item/:order_item_id", controllers.EditOrderItem())
}
