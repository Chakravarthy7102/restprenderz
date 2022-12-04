package controllers

import (
	"context"
	"golang-resto-management/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMenus() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancle = context.WithTimeout(context.Background(), 10*time.Second)

		result, err := menuCollection.Find(context.TODO(), bson.M{})

		defer cancle()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   err.Error(),
				"message": "error while getting the menues",
			})
		}
		var menus []bson.M

		if err = result.All(ctx, &menus); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusAccepted, menus)

	}
}

func GetMenu() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancle = context.WithTimeout(context.Background(), 10*time.Second)

		menu_id := c.Param("menu_id")

		var menu models.Menu

		err := menuCollection.FindOne(ctx, bson.M{"menu_id": menu_id}).Decode(&menu)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error while finding the menu",
			})
			defer cancle()
			return
		}

		defer cancle()

		c.JSON(http.StatusOK, menu)

	}
}

func CreateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)

		var menu models.Menu
		err := c.BindJSON(&menu)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			defer cancle()
			return
		}

		validationErr := validate.Struct(&menu)

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"menu": err.Error(),
			})
			defer cancle()
			return
		}
		menu.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.ID = primitive.NewObjectID()
		menu.Menu_Id = menu.ID.Hex()

		result, insertErr := menuCollection.InsertOne(ctx, menu)

		defer cancle()

		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error while inserting the document",
				"error":   insertErr.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, result)

	}
}

//checks for if the given start and end date's are valid timelines
func inTimeSpan(start, end, check time.Time) bool {
	return start.After(check) && end.After(start)
}

func EditMenu() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		var menu models.Menu
		err := c.BindJSON(&menu)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			defer cancel()
			return
		}

		menuId := c.Param("menu_id")

		filter := bson.M{"menu_id": menuId}

		var update_object primitive.D

		if menu.Start_Date != nil && menu.End_Date != nil {

			if !inTimeSpan(*menu.Start_Date, *menu.End_Date, time.Now()) {
				msg := "Please pass the valid start and end dates"

				c.JSON(http.StatusBadRequest, gin.H{
					"error": msg,
				})
				defer cancel()
				return
			}

			update_object = append(update_object, bson.E{Key: "start_date", Value: menu.Start_Date})
			update_object = append(update_object, bson.E{Key: "end_date", Value: menu.End_Date})
		}

		if menu.Name != "" {
			update_object = append(update_object, bson.E{Key: "name", Value: menu.Name})
		}

		if menu.Category != "" {
			update_object = append(update_object, bson.E{Key: "category", Value: menu.Category})
		}

		menu.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		update_object = append(update_object, bson.E{Key: "updated_at", Value: menu.Updated_At})

		upsert := true

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := menuCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{Key: "$set", Value: update_object},
			},
			&opt,
		)

		if err != nil {
			msg := "menu update failed"

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": msg,
			})

			defer cancel()
			return
		}

		defer cancel()

		c.JSON(http.StatusOK, gin.H{
			"message": "update was successfull",
			"result":  result,
		})

	}
}
