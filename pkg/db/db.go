package db

import (
	"profile_service/pkg/errors"
	"profile_service/pkg/models"
)

type UserDAO interface {
	GetByUsername(username string) *models.User
	GetById(id int) *models.User
	AddReceiver(userId, receiverId int) error
	RemoveReceiver(userId, receiverId int) error
}

var Users map[int]*models.User

func init() {
	Users = map[int]*models.User{
		1: {
			ID:        1,
			Username:  "user1",
			Receivers: []int{1},
		},
		2: {
			ID:       2,
			Username: "user2",
		},
		3: {
			ID:       3,
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

func (*InMemroyUserDAO) GetById(id int) *models.User {
	user, ok := Users[id]
	if !ok {
		return nil
	}

	return user
}

func (*InMemroyUserDAO) AddReceiver(userId, receiverId int) error {
	user, ok := Users[userId]
	if !ok {
		return errors.NewUserDAOError(errors.UserNotFoundInDB, nil)
	}
	_, ok = Users[receiverId]
	if !ok {
		return errors.NewUserDAOError(errors.ReceiverNotFoundInDB, nil)
	}
	user.Receivers = append(user.Receivers, receiverId)

	return nil
}

func (*InMemroyUserDAO) RemoveReceiver(userId, receiverId int) error {
	user, ok := Users[userId]
	if !ok {
		return errors.NewUserDAOError(errors.UserNotFoundInDB, nil)
	}
	_, ok = Users[receiverId]
	if !ok {
		return errors.NewUserDAOError(errors.ReceiverNotFoundInDB, nil)
	}

	for ind, id := range user.Receivers {
		if id == receiverId {
			user.Receivers = append(user.Receivers[:ind], user.Receivers[ind+1:]...)
			return nil
		}
	}
	return errors.NewUserDAOError(errors.ReceiverNotInList, nil)
}
