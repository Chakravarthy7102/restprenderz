package controllers

import (
	"context"
	"golang-resto-management/database"
	"golang-resto-management/helpers"
	"golang-resto-management/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

//this part just returns the handler function to the router.
//no problem even if you write direct handler itself.
func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)

		limit, err := strconv.ParseInt(c.Query("limit"), 2, 2)
		if err != nil || limit < 1 {
			limit = 1
		}
		page, err := strconv.ParseInt(c.Query("page"), 2, 2)

		if err != nil || page < 1 {
			page = 1
		}

		skip := limit * (page - 1)

		match := bson.D{
			{
				"$match", bson.D{{}},
			},
		}

		project := bson.D{
			{
				"$project", bson.D{
					{"_id", 0},
					{"total_count", 1},
					{"user_item", bson.D{
						{"$slice", []interface{}{"$data", skip, limit}},
					},
					},
				},
			},
		}

		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
			match,
			project,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "something went wrong while fetching users",
				"error":   err.Error(),
			})

			defer cancle()
		}

		defer cancle()

		// bottom_state := bson.D{
		// 	{
		// 		"$skip", skip,
		// 	},
		// 	{
		// 		"$limit", limit,
		// 	},
		// }

		// options := options.FindOptions{
		// 	Skip:  &skip,
		// 	Limit: &limit,
		// }

		// users_cursor, err := userCollection.Find(ctx, bson.D{}, &options)

		// userCollection.Aggregate(ctx, mongo.Pipeline{
		// 	match,
		// 	bottom_state,
		// })

		var users []models.User

		if err := result.All(ctx, &users); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "something went wrong while getting the user",
				"error":   err.Error(),
			})
		}

		c.JSON(http.StatusOK, users)

	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)

		userId := c.Param("userId")

		var user models.User

		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)

		defer cancle()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "something went wrong while getting the user",
				"errpr":   err.Error(),
			})
		}

		c.JSON(http.StatusOK, user)

	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {

		//convert the json data thats coming from the client to that something golang understands.

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Please pass the required fileds",
				"error":   err.Error(),
			})

			defer cancel()
			return
		}

		if err := validate.Struct(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "validation error",
				"error":   err.Error(),
			})
		}

		//find the user with the given data , check that user exists.

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "something went wrong while finding the user",
				"error":   err.Error(),
			})
		}

		//then verify the password.
		correct, message := VerifyPassword(*foundUser.Password, *user.Password)

		if !correct {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": message,
			})
			defer cancel()
			return
		}

		//generate the token for the user.

		token, refreshToken, _ := helpers.GenerateAllTokens(*foundUser.Email, *foundUser.First_Name, *foundUser.Last_name, foundUser.User_id)

		helpers.UpdateAllTokens(token, refreshToken, foundUser.User_id)

		c.JSON(http.StatusOK, foundUser)

		//return statusOK

	}
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		//convert the json data thats coming from the client to that something golang understands.

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}

		//validate the payload based on the user struct

		if validateError := validate.Struct(user); validateError != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": validateError.Error(),
			})
		}

		//check if unique or not.

		count, err := userCollection.CountDocuments(ctx, bson.M{
			"email": user.Email,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "something went wrong",
				"error":   err.Error(),
			})

			defer cancel()
			return
		}

		if count != 0 {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Email is already in use",
			})

			defer cancel()
			return
		}

		//hash the password

		password := HashPassword(*user.Password)

		user.Password = &password

		//also check whether the phone number is unique

		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "error while matching the phone number",
				"error":   err.Error(),
			})
			defer cancel()
			return
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "The number already exit's in our database",
			})
			defer cancel()
			return
		}

		//create some extra details for the user object - created_at, updated_at and ID

		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()

		//generate and create refresh token

		token, refresh_token, _ := helpers.GenerateAllTokens(*user.Email, *user.First_Name, *user.Last_name, *&user.User_id)

		user.Token = &token
		user.Refresh_Token = &refresh_token
		//if every thing is ok then insert the user document into the database.
		insertReuslt, err := userCollection.InsertOne(ctx, user)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "error while creating the error",
				"error":   err.Error(),
			})
			defer cancel()

			return

		}
		// and return the 201 to the client refering to the success.

		c.JSON(http.StatusOK, insertReuslt)

	}
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(passwordInDatabase string, userEnteredPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(passwordInDatabase), []byte(userEnteredPassword))

	if err != nil {
		return false, "Incorrect Email or Password."
	}

	return true, "Password matched"
}
