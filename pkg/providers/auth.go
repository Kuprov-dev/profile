package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"profile_service/pkg/conf"
	"profile_service/pkg/errors"
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

	var user models.User

	if err != nil {
		log.Println("Credentials marshaling error")
		return user, errors.NewRequestError(400, errors.CredsMarshalingError, err)
	}

	addr := service.config.GetAuthServiceAddr() + "/me"

	resp, err := http.Post(addr, "application/json",
		bytes.NewBuffer(credsJSON))

	if err != nil {
		log.Println("Auth service is unavailable")
		return user, errors.NewRequestError(503, errors.AuthServiceUnavailableError, err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Read response body error")
		return user, errors.NewRequestError(500, errors.ClientRequestError, err)
	}

	fmt.Println("BODY", string(body), resp.StatusCode)

	switch resp.StatusCode {
	case 400:
		return user, errors.NewRequestError(400, errors.BadRequestError, err)
	case 401:
		return user, errors.NewRequestError(401, errors.UnauthorisedError, err)
	case 403:
		return user, errors.NewRequestError(403, errors.ForbiddenError, err)
	case 200:
		err = json.Unmarshal(body, &user)
		if err != nil {
			log.Println("Unmarshal auth response error")
			return user, errors.NewRequestError(500, errors.ClientRequestError, err)
		}
		return user, nil
	default:
		return user, errors.NewRequestError(502, errors.AuthServiceBadGatewayError, err)
	}

}
