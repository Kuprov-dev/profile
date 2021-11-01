package db

import (
	"context"
	"log"
	"profile_service/pkg/conf"
	"profile_service/pkg/errors"
	"profile_service/pkg/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBUserDAO struct {
	db             *mongo.Database
	userCollection string
}

func ConnectMongoDB(ctx context.Context, config *conf.Config) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.GetDatabaseUri()))
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client.Database(config.Database.DBName), nil
}

func NewMongoDBUserDAO(ctx context.Context, config *conf.Config) *MongoDBUserDAO {
	db, err := ConnectMongoDB(ctx, config)

	if err != nil {
		log.Fatal(err)
	}

	return &MongoDBUserDAO{db: db, userCollection: "users"}
}
func (dao *MongoDBUserDAO) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	collection := dao.db.Collection(dao.userCollection)

	filter := bson.D{
		{"comments", bson.D{{"$gt", 300}}},
		{"tags", bson.D{{"$elemMatch", bson.M{"$eq": "programming"}}}},
	}

	var user models.User
	err := collection.FindOne(ctx, filter, nil).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (*MongoDBUserDAO) GetByUUID(ctx context.Context, userUUID primitive.ObjectID) (*models.User, error) {
	// user, ok := Users[userUUID]
	// if !ok {
	// 	return nil
	// }

	// return user
	return nil
}

func (*MongoDBUserDAO) AddReceiver(ctx context.Context, userUUID primitive.ObjectID, recieverEmail string) error {
	// user, ok := Users[userUUID]
	// if !ok {
	// 	return errors.NewUserDAOError(errors.UserNotFoundInDB, nil)
	// }

	// ok = utils.ContainsString(user.Receivers, recieverEmail)
	// if ok {
	// 	return errors.NewUserDAOError(errors.DublicateReceiver, nil)
	// }

	// user.Receivers = append(user.Receivers, recieverEmail)

	return nil
}

func (*MongoDBUserDAO) RemoveReceiver(ctx context.Context, userUUID primitive.ObjectID, recieverEmail string) error {
	// user, ok := Users[userUUID]
	// if !ok {
	// 	return errors.NewUserDAOError(errors.UserNotFoundInDB, nil)
	// }

	// ok = utils.ContainsString(user.Receivers, recieverEmail)
	// if !ok {
	// 	return errors.NewUserDAOError(errors.ReceiverNotFoundInDB, nil)
	// }

	// for ind, email := range user.Receivers {
	// 	if email == recieverEmail {
	// 		user.Receivers = append(user.Receivers[:ind], user.Receivers[ind+1:]...)
	// 		return nil
	// 	}
	// }
	return errors.NewUserDAOError(errors.ReceiverNotInList, nil)
}
