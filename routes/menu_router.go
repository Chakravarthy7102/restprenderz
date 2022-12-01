package routes

import (
	"github.com/gin-gonic/gin"

	"golang-resto-management/controllers"
)

func MenuRoutes(incommingRoutes *gin.Engine) {
	incommingRoutes.GET("/menu", controllers.GetMenus())
	incommingRoutes.GET("/menu/:menu_id", controllers.GetMenu())
	incommingRoutes.POST("/menu", controllers.CreateMenu())
	incommingRoutes.PATCH("/menu/:menu_id", controllers.EditMenu())
}
