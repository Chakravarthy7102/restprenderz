package routes

import (
	"github.com/gin-gonic/gin"

	"golang-resto-management/controllers"
)

func TableRoutes(incommingRoutes *gin.Engine) {
	incommingRoutes.GET("/tables", controllers.GetTables())
	incommingRoutes.GET("/table/:table_id", controllers.GetTable())
	incommingRoutes.POST("/table", controllers.CreateTable())
	incommingRoutes.PATCH("/table/:table_id", controllers.EditTable())
}
