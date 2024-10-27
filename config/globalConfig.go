package config

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	MongoC          *mongo.Client
	MongoCollection *mongo.Collection
	UserCollection  *mongo.Collection
	MongoURI        string
	Port            string
	Secrete         string
)

func LoadEnv() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error opening environment: %v", err)
	}
	Port = os.Getenv("APP_PORT")
	MongoURI = os.Getenv("MONGO_URI")
	Secrete = os.Getenv("SECRETE")
	return err
}

func InitDB() error {
	var err error
	clientOption := options.Client().ApplyURI(MongoURI)
	MongoC, err = mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		return err
	}
	err = MongoC.Ping(context.TODO(), nil)
	if err != nil {
		return err
	}

	MongoCollection = MongoC.Database("ebook-unga").Collection("books")
	UserCollection = MongoC.Database("ebook-unga").Collection("users")

	return err
}
