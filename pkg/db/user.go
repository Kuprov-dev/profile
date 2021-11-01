package db

import (
	"context"
	"profile_service/pkg/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserDAO interface {
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetByUUID(ctx context.Context, userUUID primitive.ObjectID) (*models.User, error)
	AddReceiver(ctx context.Context, userUUID primitive.ObjectID, recieverEmail string) error
	RemoveReceiver(ctx context.Context, userUUID primitive.ObjectID, receiverEmail string) error
}
