package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Food struct {
	ID         primitive.ObjectID `bson:"_id"`
	Name       string             `bson:"name"`
	Price      string             `bson:"price"`
	Food_Image string             `bson:"food_image"`
	Created_At time.Time          `bson:"created_at"`
	Updated_At time.Time          `bson:"updated_at"`
	Food_Id    string             `bson:"food_id"`
	Menu_Id    string             `bson:"menu_id"`
}
