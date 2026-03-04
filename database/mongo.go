/**

package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB is a global variable so the whole app can use the database
var DB *mongo.Database

func ConnectDB() {

	uri := "mongodb://localhost:27017"

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println("connected to mongodb")

	DB = client.Database("go_commerce")

}

**/

package database

import (
	"context"
	"fmt"
	"go-ecommerce/config"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB is global
var DB *mongo.Database

func ConnectDB() {

	// ✅ read Mongo URI from environment
	uri := config.GetEnv("MONGO_URI")
	if uri == "" {
		panic("MONGO_URI not set in environment")
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(uri).SetMaxPoolSize(50))
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println("✅ Connected to MongoDB")

	DB = client.Database("go_commerce")

	createIndexes()

}

func createIndexes() {
	userCollection := DB.Collection("users")

	indexModel := mongo.IndexModel{
		Keys:    bson.M{"email": 1},
		Options: options.Index().SetUnique(true),
	}

	userCollection.Indexes().CreateOne(context.Background(), indexModel)
}
