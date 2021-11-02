package db

import (
	"context"
	"log"
	"profile_service/pkg/conf"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoDBConnection *mongo.Database

func ConnectMongoDB(ctx context.Context, config *conf.Config) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.GetDatabaseUri()))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	mongoDBConnection = client.Database(config.Database.DBName)
	log.Println("Connection to DB success")
}

func GetMongoDBConnection() *mongo.Database {
	return mongoDBConnection
}
