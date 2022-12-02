package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderItem struct {
	ID            primitive.ObjectID `bson:"_id"`
	Quantity      *string            `json:"quantity" validate:"required,eq=S|eq=M|eq=L"`
	Unit_Price    *float64           `json:"unit_price" validate:"required"`
	Updated_At    time.Time          `json:"created_at"`
	Created_At    time.Time          `json:"updated_at"`
	Food_Id       *string            `json:"food_id" validate:"required"`
	Order_id      string             `json:"order_id" validate:"required"`
	Order_Item_id string             `json:"order_item_id" validate:"required"`
}
