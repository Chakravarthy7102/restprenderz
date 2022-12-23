package helpers

import (
	"fmt"
	"golang-resto-management/database"
	"log"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
)

type SignedDetails struct {
	Email      string
	First_Name string
	Last_Name  string
	Uid        string
	jwt.StandardClaims
}

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

var SECRET_KEY = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email string, firstaName string, lastName string, userId string) (signedToken string, signedRefreshToken string, err error) {
	cliams := &SignedDetails{
		Email:      email,
		First_Name: firstaName,
		Last_Name:  lastName,
		Uid:        userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	token, token_err := jwt.NewWithClaims(jwt.SigningMethodHS256, cliams).SigningString()

	refreshToken, refreshToken_err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SigningString()

	if token_err != nil {
		return "", "", token_err
	}

	if refreshToken_err != nil {
		return "", "", refreshToken_err
	}

	return token, refreshToken, nil
}

func UpdateAllTokens(token string, refreshToken string, userId string) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	var updated_object primitive.D

	updated_object = append(updated_object, bson.E{"token", token})
	updated_object = append(updated_object, bson.E{"refresh_token", refreshToken})

	Updated_At, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	updated_object = append(updated_object, bson.E{"updated_at", Updated_At})

	upsert := true

	filter := bson.M{"user_id": userId}

	options := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := userCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{
				"$set", updated_object,
			},
		},
		&options,
	)

	defer cancel()

	if err != nil {
		log.Panic(err)
		return
	}

	return

}
func ValidateToken(signedToken string) (claims *SignedDetails, message string) {

	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	//token is invalid

	claims, ok := token.Claims.(*SignedDetails)

	if !ok {
		message = fmt.Sprintf("The Token is invalid.")
		message = err.Error()
		return
	}

	//token expired

	if claims.ExpiresAt < time.Now().Local().Unix() {
		message = fmt.Sprintf("Token expired")
		message = err.Error()
		return
	}

	return claims, message

}
