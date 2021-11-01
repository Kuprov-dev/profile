package db

import (
	"profile_service/pkg/conf"
	"profile_service/pkg/errors"
	"profile_service/pkg/models"
	"profile_service/pkg/utils"

	"github.com/google/uuid"
)

var Users map[uuid.UUID]*models.User

type InMemroyUserDAO struct {
}

func NewInMemoryUserDAO(config *conf.Config) *InMemroyUserDAO {
	Users = config.Database.Users
	return &InMemroyUserDAO{}
}
func (*InMemroyUserDAO) GetByUsername(username string) *models.User {
	var user *models.User

	for _, u := range Users {
		if u.Username == username {
			user = u
			break
		}
	}

	return user
}

func (*InMemroyUserDAO) GetByUUID(userUUID uuid.UUID) *models.User {
	user, ok := Users[userUUID]
	if !ok {
		return nil
	}

	return user
}

func (*InMemroyUserDAO) AddReceiver(userUUID uuid.UUID, recieverEmail string) error {
	user, ok := Users[userUUID]
	if !ok {
		return errors.NewUserDAOError(errors.UserNotFoundInDB, nil)
	}

	ok = utils.ContainsString(user.Receivers, recieverEmail)
	if ok {
		return errors.NewUserDAOError(errors.DublicateReceiver, nil)
	}

	user.Receivers = append(user.Receivers, recieverEmail)

	return nil
}

func (*InMemroyUserDAO) RemoveReceiver(userUUID uuid.UUID, recieverEmail string) error {
	user, ok := Users[userUUID]
	if !ok {
		return errors.NewUserDAOError(errors.UserNotFoundInDB, nil)
	}

	ok = utils.ContainsString(user.Receivers, recieverEmail)
	if !ok {
		return errors.NewUserDAOError(errors.ReceiverNotFoundInDB, nil)
	}

	for ind, email := range user.Receivers {
		if email == recieverEmail {
			user.Receivers = append(user.Receivers[:ind], user.Receivers[ind+1:]...)
			return nil
		}
	}
	return errors.NewUserDAOError(errors.ReceiverNotInList, nil)
}
