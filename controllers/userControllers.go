package controllers

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	helper "github.com/hauchongtang/splatbackend/functions"
	"github.com/hauchongtang/splatbackend/models"
	"github.com/hauchongtang/splatbackend/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = repository.OpenCollection(repository.Client, "users")
var taskCollection *mongo.Collection = repository.OpenCollection(repository.Client, "tasks")
var validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("login or passowrd is incorrect")
		check = false
	}

	return check, msg
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the email"})
			return
		}

		password := HashPassword(*user.Password)
		user.Password = &password

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this email already exists"})
			return
		}

		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()
		token, refreshToken, _ := helper.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, user.User_id)
		user.Token = &token
		user.Refresh_token = &refreshToken

		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := insertErr
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()

		c.JSON(http.StatusOK, resultInsertionNumber)

	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		var foundUser models.User
		c.Request.Header.Add("Access-Control-Allow-Origin", "*")

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "login or passowrd is incorrect"})
			return
		}

		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()
		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		if foundUser.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
			return
		}
		token, refreshToken, _ := helper.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, foundUser.User_id)

		helper.UpdateAllTokens(token, refreshToken, foundUser.User_id)
		err = userCollection.FindOne(ctx, bson.M{"user_id": foundUser.User_id}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, foundUser)

	}
}

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		results := make([]models.User, 0)
		c.Request.Header.Add("Access-Control-Allow-Origin", "*")
		// result := models.User{}
		filter := bson.M{}
		opts := options.Find().SetSort(bson.D{{"points", -1}})
		docCursor, err := userCollection.Find(ctx, filter, opts)

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

		log.Println(&results)
		c.JSON(http.StatusOK, &results)
	}
}

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

		log.Println(&results)
		c.JSON(http.StatusOK, &results)
	}
}

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

func GetUserById() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		c.Request.Header.Add("Access-Control-Allow-Origin", "*")
		result := models.User{}
		targetId := c.Param("id")
		filter := bson.M{"user_id": targetId}

		docCursor := userCollection.FindOne(ctx, filter)
		err := docCursor.Decode(&result)

		if err != nil {
			log.Default().Print("Unable to decode object from mongodb")
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, &result)
	}
}

func GetTasksById() gin.HandlerFunc {
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
		}
		docCursor := taskCollection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate())

		c.JSON(http.StatusOK, docCursor)
	}
}

func DeleteUserById() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		c.Request.Header.Add("Access-Control-Allow-Origin", "*")
		targetId := c.Param("id")
		objectId, err := primitive.ObjectIDFromHex(targetId)

		if err != nil {
			log.Println(err)
		}

		filter := bson.M{"_id": objectId}

		_, err = userCollection.DeleteOne(ctx, filter)

		if err != nil {
			log.Println("Failed to delete from db")
			log.Println(err)
		}

		c.JSON(http.StatusOK, "Delete Success")
	}
}

func IncreasePoints() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		c.Request.Header.Add("Access-Control-Allow-Origin", "*")
		targetId := c.Param("id")
		log.Println(targetId)
		pointsToAdd := c.Query("pointstoadd")
		_id, err := primitive.ObjectIDFromHex(targetId)

		if err != nil {
			log.Println(err)
		}

		points, err := strconv.ParseInt(pointsToAdd, 0, 64)
		filter := bson.M{"_id": _id}
		update := bson.D{
			{"$inc", bson.D{{"points", points}}},
		}
		docCursor := userCollection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetUpsert(true))

		c.JSON(http.StatusOK, docCursor)
	}
}

func UpdateModuleImportLink() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		c.Request.Header.Add("Access-Control-Allow-Origin", "*")
		targetId := c.Param("id")
		log.Println(targetId)
		linkToAdd := c.Query("linktoadd")
		// if !strings.Contains(linkToAdd, "nusmods.com/timetable") {
		// 	c.JSON(http.StatusForbidden, userCollection.FindOne(ctx, bson.M{"_id": "0"}))
		// }
		_id, err := primitive.ObjectIDFromHex(targetId)

		if err != nil {
			log.Println(err)
		}

		filter := bson.M{"_id": _id}
		update := bson.D{{"$set", bson.D{{"timetable", linkToAdd}}}}
		docCursor := userCollection.FindOneAndUpdate(ctx, filter, update)

		c.JSON(http.StatusOK, docCursor)
	}
}
