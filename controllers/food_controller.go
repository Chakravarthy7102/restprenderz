package controllers

import (
	"context"
	"golang-resto-management/database"
	"golang-resto-management/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")
var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")
var validate *validator.Validate = validator.New()

func GetFoods() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		//parsing the string to integer with strconv lib
		records_per_page, err := strconv.Atoi(c.Query("records_per_page"))

		if err != nil || records_per_page < 1 {
			records_per_page = 10
		}

		page, err := strconv.Atoi(c.Query("page"))

		if err != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * records_per_page

		match_stage := bson.D{{"$match", bson.D{{}}}}
		group_stage := bson.D{{"$group", bson.D{
			{"_id", "null"},
			{"total_count", bson.D{{"$sum", "1"}}},
			{"data", bson.D{{"$push", "$$ROOT"}}},
		}}}

		project := bson.D{
			{
				"$project", bson.D{
					{"_id", 0},
					{"total_count", 1},
					{"food_items", bson.D{
						{"$slice", []interface{}{"$data", startIndex, records_per_page}},
					}},
				},
			},
		}

		cursor, err := foodCollection.Aggregate(ctx, mongo.Pipeline{
			match_stage,
			group_stage,
			project,
		})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   err.Error(),
				"message": "something went wrong when fetching data!",
			})

			defer cancel()
			return
		}

		var allFoods []models.Food
		if err = cursor.All(ctx, &allFoods); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, allFoods)
	}
}

func GetFood() gin.HandlerFunc {
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
			return
		}
		c.JSON(http.StatusAccepted, food)
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
			defer cancle()
			return
		}

		validaionErr := validate.Struct(food)
		if validaionErr != nil {
			c.JSON(http.StatusBadRequest, validaionErr.Error())
			defer cancle()
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
		var num = food.Price
		food.Price = num

		result, insertErr := menuCollection.InsertOne(ctx, food)

		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "could not create the resource!",
			})
			defer cancle()
			return
		}
		defer cancle()
		c.JSON(http.StatusCreated, result)
	}
}

func EditFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var menu models.Menu
		var food models.Food

		food_id := c.Param("food_item_id")

		if err := c.BindJSON(&food); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"error":   err.Error(),
				"message": "something went wrong in validation",
			})

			defer cancel()
			return
		}

		var updated_object primitive.D

		if food.Name != nil {
			updated_object = append(updated_object, bson.E{Key: "name", Value: food.Name})
		}

		if food.Food_Image != nil {
			updated_object = append(updated_object, bson.E{Key: "food_image", Value: food.Food_Image})
		}

		if food.Menu_Id != nil {
			updated_object = append(updated_object, bson.E{Key: "menu_id", Value: food.Menu_Id})
		}

		if food.Price != nil {
			updated_object = append(updated_object, bson.E{Key: "price", Value: food.Price})
		}

		if food.Menu_Id != nil {
			err := menuCollection.FindOne(ctx, bson.M{"menu_id": food.Menu_Id}).Decode(&menu)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Internal error while getting the menu!",
					"error":   err.Error(),
				})
				defer cancel()
				return
			}
		}

		food.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		updated_object = append(updated_object, bson.E{Key: "updated_at", Value: food.Updated_At})
		var upsert = true
		var options = options.UpdateOptions{
			Upsert: &upsert,
		}

		updated_result, err := foodCollection.UpdateOne(ctx, bson.M{"food_id": food_id}, bson.D{
			{Key: "$set", Value: updated_object},
		}, &options)

		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong!",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, updated_result)

	}
}
