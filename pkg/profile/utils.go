package profile

import (
	"context"
	"log"
	"net/http"
	"profile_service/pkg/errors"
	"profile_service/pkg/models"
	"time"
)

// Интерфейс оборачиваемой ф-ции для Retry
type Effoctor func(context.Context) error

// Реализация паттерна Retry
// Если ресурс недоступен, то пробуем подергать его еще
func Retry(effoctor Effoctor, retries int, delay time.Duration) Effoctor {
	return func(ctx context.Context) error {
		for r := 0; ; r++ {
			log.Println("...Attempt ", r)
			err := effoctor(ctx)
			if err == nil || r > retries {
				return err
			}

			// retry only if external service is unavailable
			// TODO refactor for DRY
			ve, ok := err.(*errors.RequestError)
			if ok {
				switch {
				case ve.Errors&errors.UnauthorisedError != 0:
					return err
				case ve.Errors&errors.ForbiddenError != 0:
					return err
				}
			}

			select {
			case <-time.After(delay):
			case <-ctx.Done():
				log.Println("canceling")
				return ctx.Err()
			}

		}
	}
}

func RefreshTokenHeaders(w *http.ResponseWriter, refreshedTokenCreds *models.RefreshedTokenCreds) {
	http.SetCookie(*w, &http.Cookie{
		Name:     "Access",
		Value:    refreshedTokenCreds.AccessToken,
		Path:     "/",
		Expires:  refreshedTokenCreds.AccessExpirationTime,
		HttpOnly: true,
	})
	http.SetCookie(*w, &http.Cookie{
		Name:     "Refresh",
		Value:    refreshedTokenCreds.RefreshedToken,
		Path:     "/",
		Expires:  refreshedTokenCreds.RefreshExpirationTime,
		HttpOnly: true,
	})
}
