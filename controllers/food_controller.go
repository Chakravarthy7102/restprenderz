package controllers

import (
	"context"
	"golang-resto-management/database"
	"golang-resto-management/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")
var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")
var validate *validator.Validate = validator.New()

func GetFoods() gin.HandlerFunc {
	return func(c *gin.Context) {
		//this request must be completed within 10 seconds else the request will be terminates
		//and the resources are realesed.
		var ctx, cancle = context.WithTimeout(context.Background(), 10*time.Second)
		food_id := c.Param("food_item_id")

		var food models.Food

		err := foodCollection.FindOne(ctx, bson.M{"food_id": food_id}).Decode(&food)
		defer cancle()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "something went wrong in food handler!",
			})
		}
		c.JSON(http.StatusAccepted, food)
	}
}

func GetFood() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func CreateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancle = context.WithTimeout(context.Background(), 10*time.Second)

		var menu models.Menu
		var food models.Food

		err := c.BindJSON(&food)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		validaionErr := validate.Struct(food)
		if validaionErr != nil {
			c.JSON(http.StatusBadRequest, validaionErr.Error())
			return
		}

		err = menuCollection.FindOne(ctx, bson.M{"menu_id": food.Menu_Id}).Decode(&menu)
		defer cancle()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		food.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.ID = primitive.NewObjectID()
		food.Food_Id = food.ID.Hex()
		var num = toFixed(*food.Price, 2)
		food.Price = &num

		result, insertErr := menuCollection.InsertOne(ctx, food)

		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "could not create the resource!",
			})

			return
		}
		defer cancle()
		c.JSON(http.StatusCreated, result)
	}
}

func EditFood() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
