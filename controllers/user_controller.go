package controllers

import (
	"github.com/gin-gonic/gin"
)

//this part just returns the handler function to the router.
//no problem even if you write direct handler itself.
func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func HashPassword(password string) string {
	return ""
}

func VerifyPassword(password string) bool {
	return false
}
