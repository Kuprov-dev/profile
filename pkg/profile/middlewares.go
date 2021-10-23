package profile

import (
	"context"
	"net/http"
	"profile_service/pkg/conf"
	"profile_service/pkg/db"
	"profile_service/pkg/errors"
	"profile_service/pkg/models"
	"profile_service/pkg/providers"
	"strconv"
	"strings"
)

type ContextKey string

const ContextUserDetailsKey ContextKey = "userDetails"
const ContextUserIdKey ContextKey = "userId"

// мидлварь чтобы проверить что юзер это самое
func IsAuthenticated(config *conf.Config, userDao db.UserDAO, authService providers.AuthServiceProvider) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			creds := models.UserCredentials{
				AccessToken:  r.Header.Get("Access"),
				RefreshToken: r.Header.Get("Refresh"),
			}

			userDetails, err := checkUserIsAuthenticated(r.Context(), &creds, config, authService)

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

			ctx := context.WithValue(r.Context(), ContextUserDetailsKey, userDetails)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// мидлварь для разграничения доступа к ресурсу
func IsOwner(config *conf.Config, userDao db.UserDAO, authService providers.AuthServiceProvider) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userValue := r.Context().Value(ContextUserDetailsKey)

			user, ok := userValue.(*models.UserDetails)

			if !ok {
				makeInternalServerErrorResponse(&w)
				return
			}

			urlPath := strings.Split(r.URL.Path, "/")
			userId, _ := strconv.Atoi(urlPath[len(urlPath)-2])

			if !checkIsTheSameUser(userId, user.Username, userDao) {
				makeForbiddenErrorResponse(&w)
				return
			}

			ctx := context.WithValue(r.Context(), ContextUserIdKey, userId)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
