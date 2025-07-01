package database

import (
	"context"
	"fmt"
	"log"
	"notes_project/envs"
	"time"

	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var MongoClient *mongo.Client
var RedisClient *redis.Client

func InitDatabase() error {
	env := &envs.ServerEnvs
	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s:%s",
		env.MONGO_INITDB_ROOT_USERNAME,
		env.MONGO_INITDB_ROOT_PASSWORD,
		env.MONGO_INITDB_HOST,
		env.MONGO_INITDB_PORT)
	log.Println("URL:" + mongoURI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongo, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return err
	} else {
		MongoClient = mongo
	}
	mongoErr := MongoClient.Ping(ctx, readpref.Primary())
	if mongoErr != nil {
		return mongoErr
	}
	return nil
}

func InitRedis() error {
	redisURI := fmt.Sprintf("%s:%s", envs.ServerEnvs.REDIS_HOST, envs.ServerEnvs.REDIS_PORT)

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisURI,
		Password: "",
		DB:       0,
	})

	status := RedisClient.Ping()
	if status.Val() == "PONG" {
		return nil
	} else {
		return fmt.Errorf("Ошибка при подключении к Redis: %v", status)
	}
}
