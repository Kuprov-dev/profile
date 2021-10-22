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
// и ходит за данными юзера в сервис auth
func getUserDataFromAuthService(ctx context.Context, creds *models.UserCredentials, authService providers.AuthServiceProvider, config conf.Config) (*models.UserDetails, error) {
	var user *models.UserDetails
	var err error

	effector := func(ctx context.Context) error {
		user, err = authService.GetUserData(creds)
		log.Println("Effector ", *user, err)
		if err != nil {
			return err
		}
		return nil
	}

	effectorWithRetry := Retry(effector, config.AuthServiceRetries, time.Duration(config.AuthServiceRetryDelay)*time.Millisecond)
	err = effectorWithRetry(ctx)

	return user, err
}

// Интерактор, который инкапсулирует логику работы с AuthServiceProvider
// и ходит за проверкой валидности токена в сервис auth
func checkUserIsAuthenticated(ctx context.Context, creds *models.UserCredentials, config *conf.Config, authService providers.AuthServiceProvider) (*models.UserDetails, error) {
	var err error
	var user *models.UserDetails

	effector := func(ctx context.Context) error {
		user, err = authService.CheckUserIsAuthenticated(ctx, creds)
		log.Println("Effector ", *user, err)
		if err != nil {
			return err
		}
		return nil
	}

	effectorWithRetry := Retry(effector, config.AuthServiceRetries, time.Duration(config.AuthServiceRetryDelay)*time.Millisecond)
	err = effectorWithRetry(ctx)

	return user, err
}
