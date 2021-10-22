package providers

import (
	"bytes"
	"context"
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
	GetUserData(creds *models.UserCredentials) (*models.UserDetails, error)
	CheckUserIsAuthenticated(ctx context.Context, creds *models.UserCredentials) (*models.UserDetails, error)
}

type HttpAuthServiceProvider struct {
	Config *conf.Config
}

func NewHttpAuthServiceProvider(config *conf.Config) *HttpAuthServiceProvider {
	return &HttpAuthServiceProvider{Config: config}
}

func (service *HttpAuthServiceProvider) GetUserData(creds *models.UserCredentials) (*models.UserDetails, error) {
	credsJSON, err := json.Marshal(creds)

	var user models.UserDetails

	if err != nil {
		log.Println("Credentials marshaling error")
		return &user, errors.NewRequestError(400, errors.CredsMarshalingError, err)
	}

	url := service.Config.GetAuthServiceProfileDetailsUrl()

	resp, err := http.Post(url, "application/json",
		bytes.NewBuffer(credsJSON))

	if err != nil {
		log.Println("Auth service is unavailable")
		return &user, errors.NewRequestError(503, errors.AuthServiceUnavailableError, err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Read response body error")
		return &user, errors.NewRequestError(500, errors.ClientRequestError, err)
	}

	fmt.Println("BODY", string(body), resp.StatusCode)

	// TODO refactor for DRY using errors map for example
	switch resp.StatusCode {
	case 400:
		return &user, errors.NewRequestError(400, errors.BadRequestError, err)
	case 401:
		return &user, errors.NewRequestError(401, errors.UnauthorisedError, err)
	case 403:
		return &user, errors.NewRequestError(403, errors.ForbiddenError, err)
	case 200:
		err = json.Unmarshal(body, &user)
		if err != nil {
			log.Println("Unmarshal auth response error")
			return &user, errors.NewRequestError(500, errors.ClientRequestError, err)
		}
		return &user, nil
	default:
		return &user, errors.NewRequestError(502, errors.AuthServiceBadGatewayError, err)
	}

}

func (service *HttpAuthServiceProvider) CheckUserIsAuthenticated(ctx context.Context, creds *models.UserCredentials) (*models.UserDetails, error) {
	var user models.UserDetails

	client := http.Client{}

	url := service.Config.GetAuthServiceTokenValidationUrl()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &user, errors.NewRequestError(500, errors.ClientRequestError, err)
	}

	req.Header.Set("Authorization", creds.AccessToken)
	req.Header.Set("Refresh", creds.RefreshToken)

	resp, err := client.Do(req)
	if err != nil {
		return &user, errors.NewRequestError(503, errors.AuthServiceUnavailableError, err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		log.Println("Read response body error")
		return &user, errors.NewRequestError(500, errors.ClientRequestError, err)
	}
	log.Println("Response from auth service: ", resp.Status)

	switch resp.StatusCode {
	case 400:
		return &user, errors.NewRequestError(400, errors.BadRequestError, err)
	case 401:
		return &user, errors.NewRequestError(401, errors.UnauthorisedError, err)
	case 403:
		return &user, errors.NewRequestError(403, errors.ForbiddenError, err)
	case 200:
		err = json.Unmarshal(body, &user)
		log.Println("User:", user)
		if err != nil {
			log.Println("Unmarshal auth response error")
			return &user, errors.NewRequestError(500, errors.ClientRequestError, err)
		}
		return &user, nil
	default:
		return &user, errors.NewRequestError(502, errors.AuthServiceBadGatewayError, err)
	}

}