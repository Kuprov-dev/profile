package db

import (
	"fmt"
	"profile_service/pkg/errors"
	"profile_service/pkg/models"

	"github.com/google/uuid"
)

type UserDAO interface {
	GetByUsername(username string) *models.User
	GetByUUID(userUUID uuid.UUID) *models.User
	AddReceiver(userUUID uuid.UUID, recieverEmail string) error
	RemoveReceiver(userUUID uuid.UUID, receiverEmail string) error
}

var Users map[uuid.UUID]*models.User

func init() {
	uuids := []uuid.UUID{
		uuid.New(),
		uuid.New(),
		uuid.New(),
	}
	fmt.Println(uuids)
	Users = map[uuid.UUID]*models.User{
		uuids[0]: {
			UUID:      uuids[0],
			Username:  "user1",
			Receivers: []string{"test@mail.ru"},
		},
		uuids[1]: {
			UUID:     uuids[1],
			Username: "user2",
		},
		uuids[2]: {
			UUID:     uuids[2],
			Username: "user3",
		},
	}
}

type InMemroyUserDAO struct {
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

	ok = contains(user.Receivers, recieverEmail)
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

	ok = contains(user.Receivers, recieverEmail)
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
