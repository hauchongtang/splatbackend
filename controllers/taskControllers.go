package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/cache/v9"
	"github.com/hauchongtang/splatbackend/models"
	"github.com/hauchongtang/splatbackend/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var taskCollection *mongo.Collection = repository.OpenCollection(repository.Client, "tasks")

type taskType = models.Task
type taskAddType = models.TaskResult

// GetAllActivity gdoc
// @Summary Get all task activities
// @Description Gets all tasks from the database. Represents all activities.
// @Tags task
// @Produce json
// @Security ApiKeyAuth
// @param token header string true "Authorization token"
// @Success 200 {object} []taskType
// @Failure 404 {object} errorResult
// @Router /tasks [get]
func GetAllActivity() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		results := make([]models.Task, 0)
		c.Request.Header.Add("Access-Control-Allow-Origin", "*")
		// result := models.User{}
		filter := bson.M{}
		opts := options.Find().SetSort(bson.D{{"_id", -1}})
		docCursor, err := taskCollection.Find(ctx, filter, opts)

		if err != nil {
			log.Fatal("unable to find users")
			log.Fatal(err)
			log.Fatal(docCursor)
		}

		err = docCursor.All(context.TODO(), &results)

		if err != nil {
			log.Fatal("Unable to decode list of users")
			log.Fatal(err)
			log.Fatal(docCursor.Current)
			return
		}

		c.JSON(http.StatusOK, &results)
	}
}

// GetCachedAllActivity gdoc
// @Summary Get all task activities from cache
// @Description Gets all tasks from the cache. Represents all activities.
// @Tags task
// @Produce json
// @Security ApiKeyAuth
// @param token header string true "Authorization token"
// @Success 200 {object} []taskType
// @Failure 404 {object} errorResult
// @Router /cached/tasks [get]
func GetCachedAllActivity() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		results := make([]models.Task, 0)
		c.Request.Header.Add("Access-Control-Allow-Origin", "*")
		// result := models.User{}

		redisCache.Get(ctx, "alltaskscache", &results)

		if len(results) != 0 {
			c.JSON(http.StatusOK, &results)
			log.Default().Println("Fetched from cache!")
			return
		}

		filter := bson.M{}
		opts := options.Find().SetSort(bson.D{{"_id", -1}})
		docCursor, err := taskCollection.Find(ctx, filter, opts)

		if err != nil {
			log.Fatal("unable to find tasks")
			log.Fatal(err)
			log.Fatal(docCursor)
		}

		err = docCursor.All(context.TODO(), &results)

		if err != nil {
			log.Fatal("Unable to decode list of tasks")
			log.Fatal(err)
			log.Fatal(docCursor.Current)
			return
		}

		err = redisCache.Set(&cache.Item{
			Key:   "alltaskscache",
			Value: results,
			TTL:   time.Hour * 72,
		})

		if err != nil {
			log.Default().Println("Unable to set cache!")
		}

		c.JSON(http.StatusOK, &results)
	}
}

// AddTask godoc
// @Summary Add a task
// @Description Adds task to the database
// @Tags task
// @Param data body taskAddType true "Task details"
// @Produce json
// @Security ApiKeyAuth
// @param token header string true "Authorization token"
// @Success 200 {object} taskType
// @Failure 400 {object} errorResult
// @Router /tasks [post]
func AddTask() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		var task models.Task

		c.Request.Header.Add("Access-Control-Allow-Origin", "*")

		if err := c.BindJSON(&task); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		task.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		task.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		task.ID = primitive.NewObjectID()

		resultFromInsertTask, err := taskCollection.InsertOne(ctx, task)

		if err != nil {
			msg := err
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusOK, resultFromInsertTask)
	}
}

// GetTasksByUserId gdoc
// @Summary Get all Tasks of a particular user
// @Description Gets tasks of a particular user via userId.
// @Tags task
// @Produce json
// @Param id path string true "taskId"
// @Security ApiKeyAuth
// @param token header string true "Authorization token"
// @Success 200 {object} []taskType
// @Failure 404 {object} errorResult
// @Router /tasks/{id} [get]
func GetTasksByUserId() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		c.Request.Header.Add("Access-Control-Allow-Origin", "*")
		result := make([]models.Task, 0)
		targetId := c.Param("id")
		filter := bson.M{"user_id": targetId}
		opts := options.Find().SetSort(bson.D{{"_id", -1}})
		docCursor, err := taskCollection.Find(ctx, filter, opts)

		if err != nil {
			log.Println(err)
		}

		err = docCursor.All(context.TODO(), &result)

		if err != nil {
			log.Default().Print("Unable to decode object from mongoDB")
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, &result)
	}
}

// UpdateHiddenStatus gdoc
// @Summary Sets a task to be hidden
// @Description Updates the task via provided taskId to be hidden.
// @Tags task
// @Produce json
// @Param id path string true "userId"
// @Security ApiKeyAuth
// @param token header string true "Authorization token"
// @Success 200 {object} taskType
// @Failure 404 {object} errorResult
// @Router /tasks/{id} [put]
func UpdateHiddenStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		c.Request.Header.Add("Access-Control-Allow-Origin", "*")
		targetId := c.Param("id")
		_id, err := primitive.ObjectIDFromHex(targetId)

		if err != nil {
			log.Println(err)
		}

		filter := bson.M{"_id": _id}
		update := bson.D{
			{"$set", bson.D{{"hidden", false}}},
			{"$set", bson.D{{"updated_at", time.Now()}}},
		}
		docCursor := taskCollection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate())

		c.JSON(http.StatusOK, docCursor)
	}
}

func GetMostPopularModule() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		c.Request.Header.Add("Access-Control-Allow-Origin", "*")

		docCursor, err := taskCollection.Aggregate(ctx, mongo.Pipeline{
			{{"$group", bson.D{{"count", bson.D{{"$sum", 1}}}, {"_id", bson.D{{"module_code", "$module_code"}}}}}},
		})

		if err != nil {
			log.Println(err)
		}

		var results []bson.M
		if err = docCursor.All(ctx, &results); err != nil {
			panic(err)
		}
		if err := docCursor.Close(ctx); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, results)
	}
}
