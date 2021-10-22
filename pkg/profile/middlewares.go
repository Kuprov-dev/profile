package profile

import (
	"context"
	"fmt"
	"net/http"
	"profile_service/pkg/conf"
	"profile_service/pkg/errors"
	"profile_service/pkg/models"
	"profile_service/pkg/providers"
)

type ContextKey string

const ContextUserKey ContextKey = "user"

// мидлварь чтобы проверить что юзер это самое
func IsAuthenticated(next http.HandlerFunc, config *conf.Config, authService providers.AuthServiceProvider) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Access: ", r.Header.Get("Access"))
		fmt.Println("Refresh: ", r.Header.Get("Refresh"))

		creds := models.UserCredentials{
			AccessToken:  r.Header.Get("Access"),
			RefreshToken: r.Header.Get("Refresh"),
		}

		user, err := checkUserIsAuthenticated(r.Context(), &creds, config, authService)

		if err != nil {
			ve, ok := err.(*errors.RequestError)

			if !ok {
				makeInternalServerErrorResponse(&w)
				return
			}

			// TODO refactor for DRY
			switch {
			case ve.Errors&errors.CredsMarshalingError != 0:
				makeInternalServerErrorResponse(&w)
				return
			case ve.Errors&errors.ClientRequestError != 0:
				makeInternalServerErrorResponse(&w)
				return
			case ve.Errors&errors.BadRequestError != 0:
				makeBadRequestErrorResponse(&w, "Something goes wrong.")
				return
			case ve.Errors&errors.UnauthorisedError != 0:
				makeUnathorisedErrorResponse(&w)
				return
			case ve.Errors&errors.ForbiddenError != 0:
				makeForbiddenErrorResponse(&w)
				return
			case ve.Errors&errors.AuthServiceBadGatewayError != 0:
				makeBadGatewayErrorResponse(&w)
				return
			case ve.Errors&errors.AuthServiceUnavailableError != 0:
				makeServiceUnavailableErrorResponse(&w)
				return
			default:
				makeBadRequestErrorResponse(&w, "bruuuh.")
				return
			}
		}

		ctx := context.WithValue(r.Context(), ContextUserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
