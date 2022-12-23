package middleware

import (
	"golang-resto-management/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")

		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "No Token Provided",
			})
			c.Abort()
			return
		}

		cliams, err := helpers.ValidateToken(clientToken)

		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			c.Abort()
			return
		}

		c.Set("email", cliams.Email)
		c.Set("first_name", cliams.First_Name)
		c.Set("last_name", cliams.Last_Name)
		c.Set("uid", cliams.Uid)

		c.Next()
	}
}
