package controllers

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/cache/v9"
	errors "github.com/hauchongtang/splatbackend/errors"
	helper "github.com/hauchongtang/splatbackend/functions"
	"github.com/hauchongtang/splatbackend/models"
	"github.com/hauchongtang/splatbackend/rediscache"
	"github.com/hauchongtang/splatbackend/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var userRepository *repository.UserRepository = repository.NewUserRepository(repository.Client, context.TODO())
var userCollection *mongo.Collection = repository.OpenCollection(repository.Client, "users")
var taskCollection *mongo.Collection = repository.OpenCollection(repository.Client, "tasks")
var redisCache = rediscache.Cache
var validate = validator.New()

type userType = models.User
type userSignUp = models.SignUp
type signUpResult = models.SignUpResult
type errorResult = errors.ErrorModel
type userLogin = models.LoginModel

type taskType = models.Task

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

// SignUp godoc
// @Summary User sign up
// @Description Responds with userId
// @Tags authentication
// @Param data body userSignUp true "New user credentials"
// @Produce json
// @Success 200 {object} signUpResult
// @Failure 400 {object} errorResult
// @Router /users/signup [post]
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

		// Set cache for user_id
		insertErr = redisCache.Set(&cache.Item{
			Key:   user.User_id,
			Value: user,
			TTL:   time.Hour * 72,
		})
		if insertErr != nil { // Not fatal if cache fails to set
			msg := insertErr
			log.Default().Println(msg)
		}

		resultsCache := make([]models.User, 0)
		redisCache.Get(ctx, "alluserscache", &resultsCache)
		resultsCache = append(resultsCache, user)
		insertErr = redisCache.Set(&cache.Item{
			Key:   "alluserscache",
			Value: resultsCache,
			TTL:   time.Hour * 1,
		})
		if insertErr != nil { // not fatal
			msg := insertErr
			log.Default().Println(msg)
		}

		defer cancel()

		c.JSON(http.StatusOK, resultInsertionNumber)

	}
}

// Login godoc
// @Summary User log in
// @Description Responds with user details, including OAuth2 tokens.
// @Tags authentication
// @Param data body userLogin true "Sign in credentials"
// @Produce json
// @Success 200 {object} userType
// @Failure 500 {object} errorResult
// @Router /users/login [post]
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

// GetUsers gdoc
// @Summary Get a all users
// @Description Gets all users from database directly. Use it to test whether cache is updated correctly.
// @Tags user
// @Produce json
// @Security ApiKeyAuth
// @param token header string true "Authorization token"
// @Success 200 {object} []userType
// @Failure 404 {object} errorResult
// @Router /users [get]
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

// GetCachedUsers gdoc
// @Summary Get a all users from cache
// @Description Gets all users from cache.
// @Tags user
// @Produce json
// @Security ApiKeyAuth
// @param token header string true "Authorization token"
// @Success 200 {object} []userType
// @Failure 404 {object} errorResult
// @Router /cached/users [get]
func GetCachedUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		resultsCache := make([]models.User, 0)

		redisCache.Get(ctx, "alluserscache", &resultsCache)
		if len(resultsCache) != 0 {
			c.JSON(http.StatusOK, &resultsCache)
			log.Default().Println("Fetched from cache!")
			return
		}

		results, err := userRepository.FindUsers(ctx)

		if err != nil {
			c.JSON(http.StatusNotFound, err)
			return
		}

		redisCache.Set(&cache.Item{
			Key:   "alluserscache",
			Value: results,
			TTL:   time.Hour * 1,
		})

		c.JSON(http.StatusOK, &results)
	}
}

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

		log.Println(&results)
		c.JSON(http.StatusOK, &results)
	}
}

// AddTask godoc
// @Summary Add a task
// @Description Adds task to the database
// @Tags task
// @Param data body taskType true "Task details"
// @Produce json
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

// GetUserById gdoc
// @Summary Get a User by id from database
// @Description Gets a user from database. Use this to check if the cache is updated compared to the database.
// @Tags user
// @Produce json
// @Param id path string true "userId"
// @Security ApiKeyAuth
// @param token header string true "Authorization token"
// @Success 200 {object} userType
// @Failure 404 {object} errorResult
// @Router /users/{id} [get]
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

// GetCachedUserById gdoc
// @Summary Get a User by id from cache
// @Description Gets a user from the cache if there is a hit. This is the default endpoint.
// @Tags user
// @Produce json
// @Param id path string true "userId"
// @Security ApiKeyAuth
// @param token header string true "Authorization token"
// @Success 200 {object} userType
// @Failure 404 {object} errorResult
// @Router /cached/users/{id} [get]
func GetCachedUserById() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		c.Request.Header.Add("Access-Control-Allow-Origin", "*")
		resultCache := models.User{}
		targetId := c.Param("id")
		redisCache.Get(ctx, targetId, &resultCache)
		if !resultCache.ID.IsZero() {
			fmt.Println("Result from cache!")
			c.JSON(http.StatusOK, &resultCache)
			return
		} else {
			// filter := bson.M{"user_id": targetId}
			// docCursor := userCollection.FindOne(ctx, filter)
			// err := docCursor.Decode(&result)
			result, err := userRepository.FindUserById(ctx, targetId)

			if err != nil {
				log.Default().Print("Unable to find user", targetId)
				log.Default().Println(err)
				c.JSON(http.StatusNotFound, "Unable to find user in database!")
				return
			}

			err = redisCache.Set(&cache.Item{
				Key:   targetId,
				Value: result,
				TTL:   time.Hour * 72,
			})

			if err != nil {
				log.Default().Println("Unable to set result into cache!")
			}

			c.JSON(http.StatusOK, &result)
		}
	}
}

func GetCachedUserResultById(targetId string) *models.User {
	ctx := context.Background()
	resultCache := models.User{}
	redisCache.Get(ctx, targetId, &resultCache)
	if !resultCache.ID.IsZero() {
		fmt.Println("Result from cache!")
		return &resultCache
	} else {
		// filter := bson.M{"user_id": targetId}
		// docCursor := userCollection.FindOne(ctx, filter)
		// err := docCursor.Decode(&result)
		result, err := userRepository.FindUserById(ctx, targetId)

		if err != nil {
			log.Default().Print("Unable to find user", targetId)
			log.Default().Println(err)
			return nil
		}

		err = redisCache.Set(&cache.Item{
			Key:   targetId,
			Value: result,
			TTL:   time.Hour * 72,
		})

		if err != nil {
			log.Default().Println("Unable to set result into cache!")
		}

		return result
	}
}

type queryStruct struct {
	name        string
	dataStr     string
	initialized bool
}

// ModifyParticulars gdoc
// @Summary Modify user particulars
// @Description Change user particulars
// @Tags user
// @Produce json
// @Param id path string true "userId"
// @Param first_name query string false "First name"
// @Param last_name query string false "Last name"
// @Param email query string false "Email"
// @Param password query string false "Password"
// @Security ApiKeyAuth
// @param token header string true "Authorization token"
// @Success 200 {object} userType
// @Failure 404 {object} errorResult
// @Router /users/update/{id} [put]
func ModifyParticulars() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		c.Request.Header.Add("Access-Control-Allow-Origin", "*")
		result := models.User{}
		targetId := c.Param("id")
		firstName, firstNameValid := c.GetQuery("first_name")
		lastName, lastNameValid := c.GetQuery("last_name")
		email, emailValid := c.GetQuery("email")
		pwChange, pwValid := c.GetQuery("password")
		filter := bson.M{"user_id": targetId}
		toUpdate := bson.M{}

		arr := [4]queryStruct{
			{"first_name", firstName, firstNameValid},
			{"last_name", lastName, lastNameValid},
			{"email", email, emailValid},
			{"password", pwChange, pwValid},
		}
		for _, value := range arr {
			if value.initialized && (value.name == "password") {
				toUpdate[value.name] = HashPassword(value.dataStr)
				continue
			}
			if value.initialized {
				toUpdate[value.name] = value.dataStr
			}
		}
		update := bson.M{
			"$set": toUpdate,
		}

		docCursor := userCollection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After))
		err := docCursor.Decode(&result)

		if err != nil {
			log.Default().Print("Unable to decode object from mongodb")
			log.Fatal(err)
		}

		err = redisCache.Set(&cache.Item{
			Key:   result.User_id,
			Value: result,
			TTL:   time.Hour * 72,
		})

		if err != nil { // Unable to set cache not fatal
			log.Default().Println(err)
		}

		c.JSON(http.StatusOK, &result)
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

// DeleteUserById gdoc
// @Summary Delete a user given a userId
// @Description Deletes a user via userId. Only admin access.
// @Tags user
// @Produce json
// @Param id path string true "userId"
// @Param adminId query string true "adminId"
// @Security ApiKeyAuth
// @param token header string true "Authorization token"
// @Success 200 {object} userType
// @Failure 404 {object} errorResult
// @Router /users/{id} [delete]
func DeleteUserById() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		c.Request.Header.Add("Access-Control-Allow-Origin", "*")
		targetId := c.Param("id")
		adminId := c.Query("adminId")
		objectId, err := primitive.ObjectIDFromHex(targetId)
		trueAdminId := os.Getenv("ADMIN_ID")

		if err != nil {
			log.Println(err)
		}

		if adminId != trueAdminId {
			c.JSON(http.StatusBadRequest, "Not an admin!")
			return
		}

		filter := bson.M{"_id": objectId}

		_, err = userCollection.DeleteOne(ctx, filter)

		if err != nil {
			log.Println("Failed to delete from db")
			log.Println(err)
		}

		err = redisCache.Delete(ctx, targetId)

		if err != nil { // failure to delete from cache will result in cache not being updated correctly
			log.Default().Println(err)
			log.Default().Println("Fail to delete", targetId, "from cache")
			c.JSON(http.StatusBadRequest, err)
			return
		}

		c.JSON(http.StatusOK, "Delete Success")
	}
}

// IncreasePoints gdoc
// @Summary Increase points of a user
// @Description Increase points of a user by specified amount.
// @Tags user
// @Produce json
// @Param id path string true "userId"
// @Param pointstoadd query string true "pointsToAdd"
// @Security ApiKeyAuth
// @param token header string true "Authorization token"
// @Success 200 {object} userType
// @Failure 404 {object} errorResult
// @Router /users/{id} [put]
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

		if err != nil {
			log.Println(err, "Unable to parse pointsToAdd")
		}

		filter := bson.M{"_id": _id}
		update := bson.D{
			{"$inc", bson.D{{"points", points}}},
		}
		docCursor := userCollection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetUpsert(true))

		result := models.User{}

		err = docCursor.Decode(&result)

		if err != nil {
			log.Default().Println("Unable to decode result ", err)
			c.JSON(http.StatusBadRequest, err)
			err = nil
		}

		result.Points = result.Points + int(points)

		err = redisCache.Set(&cache.Item{
			Key:   targetId,
			Value: result,
			TTL:   time.Hour * 72,
		})

		if err != nil { // Unable to update cache: Fatal
			log.Default().Println("Unable to update cache")
			c.JSON(http.StatusBadRequest, err)
			return
		}

		usersResult := make([]models.User, 0)
		err = redisCache.Get(ctx, "alluserscache", &usersResult)

		if err != nil { // Unable to update cache: Fatal
			log.Default().Println("Unable to update cache")
			c.JSON(http.StatusBadRequest, err)
			return
		}

		for i := 0; i < len(usersResult); i++ { // Linear Search O(N)
			if usersResult[i].User_id == targetId {
				usersResult[i] = result
				break
			}
		}

		err = redisCache.Set(&cache.Item{
			Key:   "alluserscache",
			Value: usersResult,
			TTL:   time.Hour * 1,
		})

		if err != nil {
			log.Default().Println("Unable to update alluserscache")
			c.JSON(http.StatusBadRequest, err)
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

// UpdateModuleImportLink gdoc
// @Summary Update the module import link of a user
// @Description Updates the module import link of the userId specified.
// @Tags user
// @Produce json
// @Param id path string true "userId"
// @Param linktoadd query string true "linkToAdd"
// @Security ApiKeyAuth
// @param token header string true "Authorization token"
// @Success 200 {object} userType
// @Failure 404 {object} errorResult
// @Router /users/modules/{id} [put]
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

		result := models.User{}

		err = docCursor.Decode(&result)

		if err != nil {
			log.Default().Println("Unable to decode result")
			log.Default().Println(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}

		result.Timetable = linkToAdd

		err = redisCache.Set(&cache.Item{
			Key:   targetId,
			Value: result,
			TTL:   time.Hour * 72,
		})

		if err != nil {
			log.Default().Println("Unable to update cache")
			log.Default().Panicln(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}

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
