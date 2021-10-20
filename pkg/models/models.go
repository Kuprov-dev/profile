package models

type User struct {
	Username string
}

type UserCredentials struct {
	AccessToken  string
	RefreshToken string
}
