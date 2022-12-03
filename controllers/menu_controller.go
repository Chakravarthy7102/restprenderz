package controllers

import (
	"context"
	"golang-resto-management/database"
	"golang-resto-management/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")

func GetMenus() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancle = context.WithTimeout(context.Background(), 10*time.Second)

		result, err := menuCollection.Find(context.TODO(), bson.M{})

		if err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{
				"error":err.Error(),
				"message":"error while getting the menues"
			})
		}
		var menus []bson.M

		if err = result.All(ctx,&menus); err != nil{
			log.Fatal(err)
		}

		c.JSON(http.StatusAccepted,menus)


		

		

	}
}

func GetMenu() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func CreateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func EditMenu() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
