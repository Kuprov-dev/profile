package db

import "profile_service/pkg/models"

type UserDAO interface {
	GetByUsername(username string) *models.User
	UpdateRefreshToken(username string, refreshToken string) *models.User
}

var Users map[string]models.User

func init() {
}

type InMemroyUserDAO struct {
}

func (*InMemroyUserDAO) GetByUsername(username string) *models.User {
	user, ok := Users[username]
	if !ok {
		return nil
	}

	return &user
}
