package repository

import (
	"context"
	"log"

	"github.com/hauchongtang/splatbackend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewUserRepository(client *mongo.Client, ctx context.Context) *UserRepository {
	err := client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	databaseName := "users"

	collection := OpenCollection(client, databaseName)

	return &UserRepository{
		collection: collection,
		ctx:        ctx,
	}
}

func (r *UserRepository) FindUserById(ctx context.Context, targetId string) (*models.User, error) {
	filter := bson.M{"user_id": targetId}
	result := models.User{}
	docCursor := r.collection.FindOne(ctx, filter)
	err := docCursor.Decode(&result)

	if err != nil {
		log.Default().Print("Unable to decode object from mongodb")
		log.Default().Print(err)
		return nil, err
	}

	return &result, nil
}

func (r *UserRepository) FindUsers(ctx context.Context) (*[]models.User, error) {
	filter := bson.M{}
	result := make([]models.User, 0)
	opts := options.Find().SetSort(bson.D{{"points", -1}})
	docCursor, err := r.collection.Find(ctx, filter, opts)

	if err != nil {
		log.Default().Println("Find all tasks failed.")
	}

	err = docCursor.All(context.TODO(), &result)

	if err != nil {
		log.Default().Print("Unable to decode object from mongodb")
		log.Default().Print(err)
		return nil, err
	}

	return &result, nil
}
