package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Food struct {
	ID         primitive.ObjectID `bson:"_id"`
	Name       string             `bson:"name" validate:"required,min=2,max=100"`
	Price      string             `bson:"price" validate:"required"`
	Food_Image string             `bson:"food_image" validate:"required"`
	Created_At time.Time          `bson:"created_at"`
	Updated_At time.Time          `bson:"updated_at"`
	Food_Id    string             `bson:"food_id"`
	Menu_Id    string             `bson:"menu_id"`
}
