package models

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// type User struct {
// 	UUID      uuid.UUID `json:"uuidd" yaml:"uuid"`
// 	Username  string    `json:"username" yaml:"username"`
// 	Receivers []string  `json:"receivers" yaml:"receivers"`
// }

type User struct {
	UUID      primitive.ObjectID `bson:"_id,omitempty"`
	Username  string             `bson:"username"`
	Receivers []string           `bson:"receivers"`
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
	Params []string `json:"params"`
}

type HTMLTeplatesListResponse struct {
	Templates []*HTMLTeplate `json:"templates"`
}
