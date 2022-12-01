package main

import (
	"golang-resto-management/database"
	"golang-resto-management/middleware"
	"golang-resto-management/routes"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var food_collection *mongo.Collection = database.OpenCollection(database.Client, "food")

func main() {
	PORT := os.Getenv("PORT")

	if PORT == "" {
		PORT = "3000"
	}

	router = gin.New()
	// logger
	router.Use(gin.Logger())

	routes.UserRoutes(router)
	router.Use(middleware.Authentication)

	routes.FoodRoutes(router)
	routes.InvoiceRoutes(router)
	routes.MenuRoutes(router)
	routes.OrderItemRoutes(router)
	routes.OrderRoutes(routes)
	routes.TableRoutes(router)

	router.Run(':' + PORT)

}
