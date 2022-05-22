package repository

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoCtx    = context.Background()
	mongoClient *mongo.Client
)

func GetMongoClient() *mongo.Client {
	err := godotenv.Load(".env")

	if err != nil {
		log.Println("error loading .env file")
	}

	mongoUri := os.Getenv("MONGODB_URI")
	if mongoClient == nil {
		mongoMinPoolSize, err := strconv.ParseUint(os.Getenv("MONGO_MIN_POOL_SIZE"), 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		mongoMaxPoolSize, err := strconv.ParseUint(os.Getenv("MONGO_MAX_POOL_SIZE"), 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		mongoMaxIdleTimeMS, err := strconv.ParseInt(os.Getenv("MONGO_MAX_IDLE_TIME_MS"), 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		clientOptions := options.Client().ApplyURI(mongoUri)

		clientOptions.SetMinPoolSize(mongoMinPoolSize)
		clientOptions.SetMaxPoolSize(mongoMaxPoolSize)
		clientOptions.SetMaxConnIdleTime(time.Duration(mongoMaxIdleTimeMS) * time.Millisecond)

		client, err := mongo.Connect(mongoCtx, clientOptions)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Connected to MongoDB!")
		return client
	}
	return mongoClient
}

//Client Database instance
var Client *mongo.Client = GetMongoClient()

//OpenCollection is a  function makes a connection with a collection in the database
func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {

	var collection *mongo.Collection = client.Database("splatbackend").Collection(collectionName)

	return collection
}
