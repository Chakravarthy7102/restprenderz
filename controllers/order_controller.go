package controllers

import (
	"context"
	"errors"
	"golang-resto-management/database"
	"golang-resto-management/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ordersCollection *mongo.Collection = database.OpenCollection(database.Client, "orders")
var tableCollection *mongo.Collection = database.OpenCollection(database.Client, "tables")

func GetOrders() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		var orders []models.Order

		cursor, err := ordersCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "internal error while fetching orders",
				"error":   err.Error(),
			})

			defer cancel()
			return
		}

		if err = cursor.All(ctx, &orders); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "internal error while fetching the orders",
				"error":   err.Error(),
			})
			defer cancel()
			return
		}

		defer cancel()

		c.JSON(http.StatusOK, orders)
	}
}

func GetOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancle = context.WithTimeout(context.Background(), 10*time.Second)
		order_id := c.Param("order_id")

		var order models.Order

		err := ordersCollection.FindOne(ctx, bson.M{"order_id": order_id}).Decode(&order)
		defer cancle()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "something went wrong in order handler!",
			})
			return
		}
		c.JSON(http.StatusAccepted, order)
	}
}

func CreateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		var order models.Order

		err := c.BindJSON(&order)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Please pass correct and required fields in the body!",
				"error":   err.Error(),
			})
			defer cancel()
			return
		}

		validation_err := validate.Struct(order)

		if validation_err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "validation error",
				"error":   err.Error(),
			})

			defer cancel()
			return
		}
		var table models.Table

		if order.Table_Id != nil {
			filter := bson.M{"table_id": order.Table_Id}
			err := tableCollection.FindOne(ctx, filter).Decode(&table)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "error while finding the table",
					"error":   err.Error(),
				})

				defer cancel()
				return
			}
		}

		order.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		order.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		order.ID = primitive.NewObjectID()
		order.Order_Id = order.ID.Hex()

		result, err := ordersCollection.InsertOne(ctx, order)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "error while creating new order",
				"error":   err.Error(),
			})
		}

		defer cancel()
		c.JSON(http.StatusOK, result)

	}
}

func EditOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		var table models.Table
		var order models.Order

		var updated_object primitive.D

		order_id := c.Param("order_id")

		if err := c.BindJSON(&order); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "something went wrong in fetching orders",
			})

			defer cancel()
			return
		}

		match := bson.M{"order_id": order_id}

		if order.Table_Id != nil {

			err := tableCollection.FindOne(ctx, bson.M{"table_id": order.Table_Id}).Decode(&table)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "something went wrong while fetching table!",
					"error":   err.Error(),
				})
				defer cancel()
				return
			}

			updated_object = append(updated_object, bson.E{Key: "table_id", Value: order.Table_Id})
		}

		var upsert = true
		var options = options.UpdateOptions{
			Upsert: &upsert,
		}

		updated_value, err := ordersCollection.UpdateOne(
			ctx,
			match,
			bson.D{
				{Key: "$set", Value: updated_object},
			},
			&options,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "something went wrong while updating the document",
				"error":   err.Error(),
			})
			defer cancel()
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, updated_value)

	}
}

func OrderItemOrderCreator(order models.Order) string {

	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)

	order.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.ID = primitive.NewObjectID()
	order.Order_Id = order.ID.Hex()

	_, err := ordersCollection.InsertOne(ctx, order)

	if err != nil {
		defer cancle()
		errors.New("Something went wrong when creating new order")
	}
	defer cancle()
	return order.Order_Id
}
