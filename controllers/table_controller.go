package controllers

import (
	"context"
	"golang-resto-management/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetTables() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		var tables []models.Table

		cursor, err := tableCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "internal error while fetching tables",
				"error":   err.Error(),
			})

			defer cancel()
			return
		}

		if err = cursor.All(ctx, &tables); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "internal error while fetching the tables",
				"error":   err.Error(),
			})
			defer cancel()
			return
		}

		defer cancel()

		c.JSON(http.StatusOK, tables)
	}
}

func GetTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancle = context.WithTimeout(context.Background(), 10*time.Second)
		table_id := c.Param("table_id")

		var table models.Table

		err := tableCollection.FindOne(ctx, bson.M{"table_id": table_id}).Decode(&table)
		defer cancle()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "something went wrong in table handler!",
			})
			return
		}
		c.JSON(http.StatusAccepted, table)
	}
}

func CreateTable() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		var table models.Table

		if err := c.BindJSON(&table); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "something went wrong while binding json",
				"error":   err.Error(),
			})
			defer cancel()
			return
		}

		if validation_error := validate.Struct(&table); validation_error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "something went wrong while validation error!",
				"error":   validation_error.Error(),
			})
			defer cancel()
			return
		}

		table.ID = primitive.NewObjectID()
		table.Table_Id = table.ID.Hex()
		table.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		table.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		inserted_item, err := tableCollection.InsertOne(ctx, table)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "something went wrong while creating table!",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, inserted_item)
	}
}

func EditTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		var table models.Table
		table_id := c.Param("table_id")

		if err := c.BindJSON(&table); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "something went wrong while binding json",
				"error":   err.Error(),
			})
			defer cancel()
			return
		}

		var updated_table primitive.D

		if table.Number_Of_guests != nil {
			updated_table = append(updated_table, primitive.E{"number_of_guests", table.Number_Of_guests})
		}

		table.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		var filter = bson.M{
			"table_id": table_id,
		}

		var upsert = true

		var options = options.UpdateOptions{
			Upsert: &upsert,
		}

		_, err := tableCollection.UpdateOne(ctx, filter, bson.D{
			{"$set", updated_table},
		}, &options)

		defer cancel()

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "something went wrong while updating the table",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "table updated"})
	}
}
