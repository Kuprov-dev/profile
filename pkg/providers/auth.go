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
	GetUserData(creds *models.UserCredentials) (*models.UserAuthDetails, error)
	CheckUserIsAuthenticated(ctx context.Context, creds *models.UserCredentials) (*models.UserAuthDetails, *models.RefreshedTokenCreds, error)
}

type HttpAuthServiceProvider struct {
	Config *conf.Config
}

func NewHttpAuthServiceProvider(config *conf.Config) *HttpAuthServiceProvider {
	return &HttpAuthServiceProvider{Config: config}
}

func (service *HttpAuthServiceProvider) GetUserData(creds *models.UserCredentials) (*models.UserAuthDetails, error) {
	credsJSON, err := json.Marshal(creds)

	var userAuthDetails models.UserAuthDetails

	if err != nil {
		log.Println("Credentials marshaling error")
		return &userAuthDetails, errors.NewRequestError(400, errors.CredsMarshalingError, fmt.Errorf("Credentials marshaling error"))
	}

	url := service.Config.GetAuthServiceProfileDetailsUrl()

	resp, err := http.Post(url, "application/json",
		bytes.NewBuffer(credsJSON))

	if err != nil {
		log.Println("Auth service is unavailable")
		return &userAuthDetails, errors.NewRequestError(503, errors.AuthServiceUnavailableError, fmt.Errorf("Auth service is unavailable"))
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Read response body error")
		return &userAuthDetails, errors.NewRequestError(500, errors.ClientRequestError, fmt.Errorf("Read response body error"))
	}

	// TODO refactor for DRY using errors map for example
	switch resp.StatusCode {
	case 400:
		return &userAuthDetails, errors.NewRequestError(400, errors.BadRequestError, err)
	case 401:
		return &userAuthDetails, errors.NewRequestError(401, errors.UnauthorisedError, err)
	case 403:
		return &userAuthDetails, errors.NewRequestError(403, errors.ForbiddenError, err)
	case 200:
		err = json.Unmarshal(body, &userAuthDetails)
		if err != nil {
			log.Println("Unmarshal auth response error")
			return &userAuthDetails, errors.NewRequestError(500, errors.ClientRequestError, err)
		}
		return &userAuthDetails, nil
	default:
		return &userAuthDetails, errors.NewRequestError(502, errors.AuthServiceBadGatewayError, err)
	}

}

func (service *HttpAuthServiceProvider) CheckUserIsAuthenticated(ctx context.Context, creds *models.UserCredentials) (*models.UserAuthDetails, *models.RefreshedTokenCreds, error) {
	var userAuthDetails models.UserAuthDetails
	var refreshedTokens *models.RefreshedTokenCreds

	client := http.Client{}

	url := service.Config.GetAuthServiceTokenValidationUrl()

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(models.TokenCredentials{
		AccessToken:  creds.AccessToken,
		RefreshToken: creds.RefreshToken,
	})

	if err != nil {
		return &userAuthDetails, refreshedTokens, errors.NewRequestError(500, errors.ClientRequestError, err)
	}

	req, err := http.NewRequest("POST", url, b)
	if err != nil {
		return &userAuthDetails, refreshedTokens, errors.NewRequestError(500, errors.ClientRequestError, err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return &userAuthDetails, refreshedTokens, errors.NewRequestError(503, errors.AuthServiceUnavailableError, err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		log.Println("Read response body error")
		return &userAuthDetails, refreshedTokens, errors.NewRequestError(500, errors.ClientRequestError, err)
	}
	log.Println("Response from auth service: ", resp.Status)

	switch resp.StatusCode {
	case 400:
		return &userAuthDetails, refreshedTokens, errors.NewRequestError(400, errors.BadRequestError, err)
	case 401:
		return &userAuthDetails, refreshedTokens, errors.NewRequestError(401, errors.UnauthorisedError, err)
	case 403:
		return &userAuthDetails, refreshedTokens, errors.NewRequestError(403, errors.ForbiddenError, err)
	case 200:
		err = json.Unmarshal(body, &userAuthDetails)
		log.Println("UserAuthDetails:", userAuthDetails)
		if err != nil {
			log.Println("Unmarshal auth response error")
			return &userAuthDetails, refreshedTokens, errors.NewRequestError(500, errors.ClientRequestError, err)
		}
		refreshedTokens = getTokenCookiesFromResponse(resp)
		return &userAuthDetails, refreshedTokens, nil
	default:
		return &userAuthDetails, refreshedTokens, errors.NewRequestError(502, errors.AuthServiceBadGatewayError, err)
	}

}
