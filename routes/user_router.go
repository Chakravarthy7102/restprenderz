package routes

import (
	"golang-resto-management/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incommingRoutes *gin.Engine) {
	incommingRoutes.GET("/users", controllers.GetUser)
	incommingRoutes.GET("/user/:userId", controllers.GetUser)
	incommingRoutes.POST("/user/signup", controllers.SignUp)
	incommingRoutes.POST("/users/login", controllers.Login)
}
