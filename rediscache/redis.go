package rediscache

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/cache/v9"
	"github.com/go-redis/redis/v9"
	"github.com/joho/godotenv"
)

var (
	redisCtx    = context.Background()
	redisClient *redis.Client
	redisCache  *cache.Cache
)

func getRedisClient() *redis.Client {
	fmt.Println("Connecting to Railway Redis Database...")

	err := godotenv.Load(".env")

	if err != nil {
		log.Println("Unable to load .env file!")
	}

	redis_uri := os.Getenv("REDIS_URI")

	if len(redis_uri) == 0 {
		log.Fatalln("No URI provided for redis client, check .env!")
	}

	opt, err := redis.ParseURL(redis_uri)

	client := redis.NewClient(opt)

	if err != nil {
		fmt.Println(err)
	}

	response, err := client.Ping(context.Background()).Result()

	if err == nil {
		fmt.Println("Connected to Railway Redis! Respnse:", response)
	} else {
		fmt.Println("Unable to connect to Railway Redis! Response:", response, err)
	}

	return client
}

func getRedisCache() *cache.Cache {
	taskcache := cache.New(&cache.Options{
		Redis: redisClient,
	})

	return taskcache
}

// Redis client database instance
var Client *redis.Client = getRedisClient()
var Cache *cache.Cache = getRedisCache()
