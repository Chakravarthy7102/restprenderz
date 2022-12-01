package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func db_instance() *mongo.Client {
	mongo_uri := "mongodb://localhost:27017"
	fmt.Println("Mongo:", mongo_uri)

	client, err := mongo.NewClient(options.Client().ApplyURI(mongo_uri))

	if err != nil {
		log.Fatal(err)
	}
	ctx1 := context.Context(context.Background())

	ctx, cancel := context.WithTimeout(ctx1, 10000) //ten seconds

	defer cancel()

	err = client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to database!")

	return client
}

var Client *mongo.Client = db_instance()

func OpenCollection(client *mongo.Client, collection_name string) *mongo.Collection {
	/**
	* var collection *mongo.Client = client.Database("resto").Collection(collection_name)
	 */
	collection := client.Database("resto").Collection(collection_name)
	return collection
}
