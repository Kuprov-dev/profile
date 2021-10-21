package models

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
