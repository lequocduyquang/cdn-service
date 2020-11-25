package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	// Client instance be exported
	Client *mongo.Client
)

// InitiateMongoClient init instance
func InitiateMongoClient() {
	var err error
	opts := options.Client().ApplyURI(os.Getenv("MONGO_URI"))
	opts.SetMaxPoolSize(5)
	if Client, err = mongo.Connect(context.TODO(), opts); err != nil {
		fmt.Printf("\nCannot connect to db: %v", err.Error())
	}

	err = Client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB successfully")
}
