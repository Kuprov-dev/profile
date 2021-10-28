package models

import (
	"html/template"
	"time"

	"github.com/google/uuid"
)

type User struct {
	UUID      uuid.UUID `json:"uuidd"`
	Username  string    `json:"username"`
	Receivers []string  `json:"receivers"`
}

type UserAuthDetails struct {
	Username string `json:"username"`
}

type UserCredentials struct {
	AccessToken  string `json:"access"`
	RefreshToken string `json:"refresh"`
}

type UserRecievers struct {
	Receivers []string `json:"receivers"`
}

type UserAddReceiver struct {
	ReceiverEmail string `json:"receiver_email"`
}

type UserRemoveReciever struct {
	ReceiverEmail string `json:"receiver_email"`
}

type HTMLTeplateCreateSchema struct {
	Name     string `json:"name"`
	Template string `json:"template"`
}

type HTMLTeplate struct {
	UUID     uuid.UUID          `json:"uuid"`
	Name     string             `json:"name"`
	Template *template.Template `json:"template"`
}

type HTMLTeplateDumpSchema struct {
	UUID     uuid.UUID `json:"uuid"`
	Name     string    `json:"name"`
	Template string    `json:"template"`
	Params   []string  `json:"params"`
}

type TokenCredentials struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshedTokenCreds struct {
	AccessToken           string
	RefreshedToken        string
	RefreshExpirationTime time.Time
	AccessExpirationTime  time.Time
}

type HTMLTeplateParsedParamsResponse struct {
	Params []string
}
