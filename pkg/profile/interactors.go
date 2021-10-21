package profile

import (
	"context"
	"log"
	"profile_service/pkg/conf"
	"profile_service/pkg/models"
	"profile_service/pkg/providers"
	"time"
)

// Интерактор, который инкапсулирует логику работы с AuthServiceProvider
func getUserDataFromAuthService(ctx context.Context, creds *models.UserCredentials, authService providers.AuthServiceProvider, config conf.Config) (*models.User, error) {
	var user models.User
	var err error

	effector := func(ctx context.Context) error {
		user, err = authService.GetUserData(creds)
		log.Println("Effector ", user, err)
		if err != nil {
			return err
		}
		return nil
	}

	effectorWithRetry := Retry(effector, config.AuthServiceRetries, time.Duration(config.AuthServiceRetryDelay)*time.Millisecond)
	err = effectorWithRetry(ctx)

	return &user, err
}
