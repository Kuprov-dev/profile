package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"profile_service/pkg/conf"
	"profile_service/pkg/models"
)

// Интерфейс для реализации паттерна Provider
// Инкапсулирует обращения к сервису auth
type AuthServiceProvider interface {
	GetUserData(creds *models.UserCredentials) (models.User, error)
}

type HttpAuthServiceProvider struct {
	config *conf.Config
}

func NewHttpAuthServiceProvider(config *conf.Config) *HttpAuthServiceProvider {
	return &HttpAuthServiceProvider{config: config}
}

func (service *HttpAuthServiceProvider) GetUserData(creds *models.UserCredentials) (models.User, error) {
	credsJSON, err := json.Marshal(creds)

	if err != nil {
		log.Println("Credentials marshaling error")
		return models.User{}, err
	}

	addr := service.config.GetAuthServiceAddr() + "/me"

	resp, err := http.Post(addr, "application/json",
		bytes.NewBuffer(credsJSON))

	if err != nil {
		log.Println("Get response from auth service error")
		return models.User{}, err
	}

	fmt.Println(resp)

	return models.User{}, nil
}
