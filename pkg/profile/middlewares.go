package profile

import (
	"context"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"profile_service/pkg/conf"
	"profile_service/pkg/db"
	"profile_service/pkg/errors"
	"profile_service/pkg/models"
	"profile_service/pkg/providers"
	"strings"
)

type ContextKey string

const ContextUserDetailsKey ContextKey = "userDetails"
const ContextUserUUIDKey ContextKey = "userUUID"

// мидлварь чтобы проверить что юзер это самое или обновить токены
func IsAuthenticatedOrRefreshTokens(config *conf.Config, authService providers.AuthServiceProvider) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			creds := models.UserCredentials{
				AccessToken:  r.Header.Get("Access"),
				RefreshToken: r.Header.Get("Refresh"),
			}

			userDetails, refreshedTokens, err := checkUserIsAuthenticated(r.Context(), &creds, config, authService)

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
			if refreshedTokens != nil {
				log.Println("Refreshing tokens")
				RefreshTokenHeaders(&w, refreshedTokens)
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

			userAuthDetails, ok := userValue.(*models.UserAuthDetails)

			if !ok {
				makeInternalServerErrorResponse(&w)
				return
			}

			urlPath := strings.Split(r.URL.Path, "/")
			// userUUID, err := uuid.Parse(urlPath[len(urlPath)-2])
			urlUUIDStr := urlPath[len(urlPath)-2]
			userUUID, err := primitive.ObjectIDFromHex(urlUUIDStr)

			if err != nil {
				makeBadRequestErrorResponse(&w, "Not valid uuid path param.")
				return
			}

			if !checkIsTheSameUser(r.Context(), userUUID, userAuthDetails.Username, userDao) {
				makeForbiddenErrorResponse(&w)
				return
			}

			ctx := context.WithValue(r.Context(), ContextUserUUIDKey, userUUID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
