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
type popularModule = struct{}

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
// @Summary Get latest task activities from cache
// @Description Gets tasks from the cache. Only the most recent 10 activities are fetched.
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

		errorMsg := redisCache.Get(ctx, "alltaskscache", &results)

		if errorMsg != nil {
			log.Default().Println(errorMsg, "Cache fetch error")
		}

		if len(results) != 0 {
			c.JSON(http.StatusOK, &results)
			log.Default().Println("Fetched from cache!")
			return
		}

		filter := bson.M{}
		opts := options.Find().SetSort(bson.D{{"_id", -1}}).SetLimit(10)
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
			TTL:   time.Hour * 1,
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
// @Failure 500 {object} errorResult
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

		userTasks := GetTasksByUserIdResult(task.User_id)

		err = redisCache.Set(&cache.Item{
			Key:   "taskOf" + task.User_id,
			Value: userTasks,
			TTL:   time.Minute * 15,
		})

		if err != nil {
			msg := err
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		// Flush alltaskscache
		err = redisCache.Delete(ctx, "alltaskscache")

		if err != nil {
			log.Fatalln(err, "Failed to flush cache")
		}

		c.JSON(http.StatusOK, resultFromInsertTask)
	}
}

// GetTasksByUserId gdoc
// @Summary Get all Tasks of a particular user
// @Description Gets tasks of a particular user via userId.
// @Tags task
// @Produce json
// @Param id path string true "userId"
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

func GetTasksByUserIdResult(targetId string) []models.Task {
	ctx := context.Background()
	result := make([]models.Task, 0)

	filter := bson.M{"user_id": targetId}
	opts := options.Find().SetSort(bson.D{{"_id", -1}})
	docCursor, err := taskCollection.Find(ctx, filter, opts)

	if err != nil {
		log.Println(err)
	}

	err = docCursor.All(context.TODO(), &result)

	if err != nil {
		log.Default().Print(err, "Unable to decode object from mongoDB")
		return result
	}

	return result
}

// GetCachedTasksByUserId gdoc
// @Summary Get latest Tasks of a particular user
// @Description Gets the latest 6 tasks of a particular user via userId.
// @Tags task
// @Produce json
// @Param id path string true "userId"
// @Security ApiKeyAuth
// @param token header string true "Authorization token"
// @Success 200 {object} []taskType
// @Failure 404 {object} errorResult
// @Router /cached/tasks/{id} [get]
func GetCachedTasksByUserId() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		c.Request.Header.Add("Access-Control-Allow-Origin", "*")
		result := make([]models.Task, 0)
		targetId := c.Param("id")

		errorMsg := redisCache.Get(ctx, "taskOf"+targetId, &result)

		if errorMsg != nil {
			log.Default().Println("Unable to get and parse from cache")
		}

		if len(result) != 0 {
			log.Default().Println("Fetched from cache!")
			c.JSON(http.StatusOK, result)
			return
		}

		filter := bson.M{"user_id": targetId}
		opts := options.Find().SetSort(bson.D{{"_id", -1}}).SetLimit(6)
		docCursor, err := taskCollection.Find(ctx, filter, opts)

		if err != nil {
			log.Println(err)
		}

		err = docCursor.All(context.TODO(), &result)

		if err != nil {
			log.Default().Print("Unable to decode object from mongoDB")
			c.JSON(http.StatusNotFound, err)
		}

		err = redisCache.Set(&cache.Item{
			Key:   "taskOf" + targetId,
			Value: result,
			TTL:   time.Minute * 15, // Prevent stale data in cache
		})

		if err != nil {
			log.Default().Println("unable to set cache")
		}

		c.JSON(http.StatusOK, &result)
	}
}

// UpdateHiddenStatus gdoc
// @Summary Sets a task to un-hide
// @Description Updates the task via provided taskId to un-hide.
// @Tags task
// @Produce json
// @Param id path string true "taskId"
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
		result := models.Task{}
		userTasks := make([]models.Task, 0)

		if err != nil {
			log.Println(err)
		}

		result.Hidden = false
		err = redisCache.Set(&cache.Item{
			Key:   "task" + targetId,
			Value: result,
			TTL:   time.Hour * 72,
		})

		if err != nil {
			log.Default().Println("Unable to set cache")
		}

		filter := bson.M{"_id": _id}
		update := bson.D{
			{"$set", bson.D{{"hidden", false}}},
			{"$set", bson.D{{"updated_at", time.Now()}}},
		}
		docCursor := taskCollection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate())

		docCursor.Decode(&result)
		result.Hidden = false
		// Update cache for user-tasks
		err = redisCache.Get(ctx, "taskOf"+result.User_id, &userTasks)

		if err != nil {
			log.Default().Println(err, "Unable to get user tasks from cache.")
			c.JSON(http.StatusOK, &result) // terminate early
			return
		}

		for i := 0; i < len(userTasks); i++ {
			if userTasks[i].ID.Hex() == targetId {
				userTasks[i] = result
			}
		}

		err = redisCache.Set(&cache.Item{
			Key:   "taskOf" + result.User_id,
			Value: userTasks,
			TTL:   time.Minute * 15,
		})

		if err != nil {
			log.Default().Println(err, "Unable to set cache")
		}

		// Flush alltaskscache
		err = redisCache.Delete(ctx, "alltaskscache")

		if err != nil {
			log.Fatalln(err, "Failed to flush cache")
		}

		c.JSON(http.StatusOK, &result)
	}
}

// GetMostPopularModule gdoc
// @Summary Get the most popular modules
// @Description Gets the module that has the most tasks done on it.
// @Tags stats
// @Produce json
// @Security ApiKeyAuth
// @param token header string true "Authorization token"
// @Success 200 {object} []popularModule
// @Failure 404 {object} errorResult
// @Router /stats/mostpopular [get]
func GetMostPopularModule() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		c.Request.Header.Add("Access-Control-Allow-Origin", "*")

		var results []bson.M
		errMsg := redisCache.Get(ctx, "mostpopularmodulescache", &results)

		if errMsg != nil {
			log.Default().Println(errMsg, "Faild to retrive from cache")
		}

		if len(results) != 0 {
			log.Default().Println("Fetched from cache!")
			c.JSON(http.StatusOK, results)
			return
		}

		docCursor, err := taskCollection.Aggregate(ctx, mongo.Pipeline{
			{{"$group", bson.D{{"count", bson.D{{"$sum", 1}}}, {"_id", bson.D{{"module_code", "$module_code"}}}}}},
		})

		if err != nil {
			log.Println(err)
		}

		if err = docCursor.All(ctx, &results); err != nil {
			panic(err)
		}
		if err := docCursor.Close(ctx); err != nil {
			panic(err)
		}

		err = redisCache.Set(&cache.Item{
			Key:   "mostpopularmodulescache",
			Value: results,
			TTL:   time.Hour * 72,
		})

		if err != nil {
			log.Default().Println("Unable to set cache")
		}

		c.JSON(http.StatusOK, results)
	}
}
