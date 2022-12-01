package routes

import (
	"golang-resto-management/controllers"

	"github.com/gin-gonic/gin"
)

func FoodRoutes(incommingRoutes *gin.Engine) {
	incommingRoutes.GET("/foods", controllers.GetFoods)
	incommingRoutes.GET("/food/:food_item_id", controllers.GetFood)
	incommingRoutes.POST("/food", controllers.CreateFood)
	incommingRoutes.PATCH("/food/:food_item_id", controllers.EditFood)
}
