package profile

import (
	"fmt"
	"net/http"
	"profile_service/pkg/conf"
	"profile_service/pkg/models"
	"profile_service/pkg/providers"
)

// Базовая ручка, чтобы ходить на auth_service/me
func ProfileDetailsHandler(config *conf.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		creds := &models.UserCredentials{}
		creds.AccessToken = r.Header.Get("Access")
		creds.RefreshToken = r.Header.Get("Refresh")

		authService := providers.NewHttpAuthServiceProvider(config)
		userData, err := getUserDataFromAuthService(r.Context(), creds, authService, *config)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("wtfff"))
		}
		fmt.Println(userData)
	}
}
