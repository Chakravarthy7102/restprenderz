package controllers

import (
	"context"
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

type OrderItemPack struct {
	Table_Id    *string
	Order_Items []models.OrderItem
}

var orderItemsCollection *mongo.Collection = database.OpenCollection(database.Client, "orderItems")

func GetOrderItems() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		var ordersItems []models.OrderItem
		cursor, err := orderItemsCollection.Find(ctx, bson.M{})

		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   err.Error(),
				"message": "internal server error",
			})
			return
		}

		if err := cursor.All(ctx, &ordersItems); err != nil {
			cursor.Close(ctx)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   err.Error(),
				"message": "internal server error",
			})
			return
		}

		c.JSON(http.StatusOK, ordersItems)

	}
}

func GetOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {

		orderId := c.Param("orderId")

		if orderId == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Please pass the orderId",
			})
		}

		allOrderedItems, err := ItemsByOrder(orderId)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   err.Error(),
				"message": "something went wrong!!",
			})
		}

		c.JSON(http.StatusOK, allOrderedItems)
	}
}

func ItemsByOrder(orderId string) ([]primitive.D, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	match_stage := bson.D{
		{
			"$match", bson.D{
				{"order_id", orderId},
			},
		},
	}

	lookup_food_stage := bson.D{
		{
			"$lookup", bson.D{
				{"from", "food"},
				{"localField", "food_id"},
				{"foreignFeild", "food_id"},
				{"as", "order"},
			},
		},
	}

	unwind_food_stage := bson.D{
		{
			"$unwind", bson.D{
				{"path", "$food"},
				{"preserveNullAndEmptyArrays", true},
			},
		},
	}

	lookup_order_stage := bson.D{
		{
			"$lookup", bson.D{
				{"from", "order"},
				{"localField", "order_id"},
				{"foreignFeild", "order_id"},
				{"as", "order"},
			},
		},
	}

	unwind_orders_stage := bson.D{
		{
			"$unwind", bson.D{
				{"path", "$order"},
				{"preserveNullAndEmptyArrays", true},
			},
		},
	}

	lookup_table_stage := bson.D{
		{
			"$lookup", bson.D{
				{"from", "table"},
				{"localField", "order.table_id"},
				{"foreignFeild", "table_id"},
				{"as", "table"},
			},
		},
	}

	unwind_table_stage := bson.D{
		{
			"$unwind", bson.D{
				{"path", "$table"},
				{"preserveNullAndEmptyArrays", true},
			},
		},
	}

	project_stage := bson.D{
		{
			"$project", bson.D{
				{"id", 0},
				{"amount", 1},
				{"total_count", 1},
				{"food_name", "$food.name"},
				{"food_image", "$food.food_image"},
				{"table_number", "$table.table_number"},
				{"table_id", "$table.table_id"},
				{"order_id", "$order.order_id"},
				{"price", "$food.price"},
				{"quantity", 1},
			},
		},
	}

	group_stage := bson.D{
		{
			"$group", bson.D{
				{
					"_id", bson.D{
						{"order_id", "$order_id"},
						{"table_id", "$table_id"},
						{"table_number", "$table_number"},
					},
				},
				{"payment_due", bson.D{{"$sum", "amount"}}},
				{"total_count", bson.D{{"$sum", 1}}},
				{"order_items", bson.D{{"$push", "$$ROOT"}}},
			},
		},
	}

	project_stage_2 := bson.D{
		{
			"$project", bson.D{
				{"id", 0},
				{"payment_due", 1},
				{"total_count", 1},
				{"order_items", 1},
				{"table_number", "$_id.table_number"},
			},
		},
	}

	cursor, err := orderItemsCollection.Aggregate(ctx, mongo.Pipeline{
		match_stage,
		lookup_food_stage,
		unwind_food_stage,
		lookup_order_stage,
		unwind_orders_stage,
		lookup_table_stage,
		unwind_table_stage,
		project_stage,
		group_stage,
		project_stage_2,
	})

	var orderItems []primitive.D

	if err != nil {
		panic(err)
	}

	if err = cursor.All(ctx, &orderItems); err != nil {
		panic(err)
	}

	defer cancel()

	return orderItems, err

}

func GetOrderItemsByOrder() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		orderItemId := c.Param("order_item_id")

		if orderItemId == "" {
			c.JSON(http.StatusOK, gin.H{
				"message": "please provide order item id",
			})
			defer cancel()
			return
		}

		var orderItem *models.OrderItem

		if err := orderItemsCollection.FindOne(ctx, bson.M{"orderItem_id": orderItemId}).Decode(&orderItem); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "something just went wrong!!",
				"error":   err.Error(),
			})

			defer cancel()
			return
		}

		defer cancel()
		c.JSON(http.StatusOK, orderItem)
	}
}

func CreateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)

		var orderItemPack OrderItemPack
		var order models.Order
		err := c.BindJSON(&orderItemPack)

		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"message": "something went wrong",
				"error":   err.Error(),
			})

			defer cancle()
			return
		}

		if validation_error := validate.Struct(&orderItemPack); validation_error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "something went wrong",
				"error":   validation_error.Error(),
			})

			defer cancle()
			return
		}

		order.Order_Date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		// var orderItemsToBeInserted []interface{}
		orderItemsToBeInserted := []interface{}{}

		order.Table_Id = orderItemPack.Table_Id
		order_id := OrderItemOrderCreator(order)

		for _, orderItem := range orderItemPack.Order_Items {
			orderItem.Order_id = order_id

			if validatation_err := validate.Struct(orderItem); validatation_err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "something really went wrong!",
					"error":   validatation_err.Error(),
				})
			}

			orderItem.ID = primitive.NewObjectID()
			orderItem.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.Order_Item_id = orderItem.ID.Hex()
			var price = *orderItem.Unit_Price
			orderItem.Unit_Price = &price
			orderItemsToBeInserted = append(orderItemsToBeInserted, orderItem)
		}

		insertedItems, err := orderItemsCollection.InsertMany(ctx, orderItemsToBeInserted)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "error while creating the items",
				"error":   err.Error(),
			})
			defer cancle()
			return
		}
		defer cancle()
		c.JSON(http.StatusCreated, insertedItems)
	}
}

func EditOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		orderItemId := c.Param("order_item_id")

		var orderedItem models.OrderItem

		if orderItemId == "" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Please pass the orderItem Id",
			})
			defer cancel()
			return
		}

		if err := c.BindJSON(&orderedItem); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "something went wrong in binding json",
				"error":   err.Error(),
			})
			defer cancel()
			return
		}

		var updated_ordered_item primitive.D

		if orderedItem.Quantity != nil {
			updated_ordered_item = append(updated_ordered_item, bson.E{"quantity", orderedItem.Quantity})
		}

		if orderedItem.Unit_Price != nil {
			updated_ordered_item = append(updated_ordered_item, bson.E{"unit_price", orderedItem.Unit_Price})
		}

		orderedItem.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updated_ordered_item = append(updated_ordered_item, primitive.E{"update_at", orderedItem.Updated_At})

		match := bson.M{
			"order_item_id": orderItemId,
		}

		upsert := true
		opts := options.UpdateOptions{
			Upsert: &upsert,
		}

		updated_results, err := orderItemsCollection.UpdateOne(ctx, match, bson.D{
			{"$set", updated_ordered_item},
		}, &opts)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "something went wrong in updating the collection.",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, updated_results)
	}
}
