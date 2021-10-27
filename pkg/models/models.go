package models

import (
	"text/template"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Receivers []int  `json:"receivers"`
}

type UserDetails struct {
	Username string `json:"username"`
}

type UserCredentials struct {
	AccessToken  string `json:"access"`
	RefreshToken string `json:"refresh"`
}

type UserRecievers struct {
	Receivers []int `json:"receivers"`
}

type UserAddReceiver struct {
	ReceiverUsername string `json:"receiver_username"`
}

type UserRemoveReciever struct {
	ReceiverUsername string `json:"receiver_username"`
}

type HTMLTeplate struct {
	Uuid     uuid.UUID         `json:"uuid"`
	Name     string            `json:"name"`
	Template template.Template `json:"template"`
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
