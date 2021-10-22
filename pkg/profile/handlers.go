package profile

import (
	"encoding/json"
	"net/http"
	"profile_service/pkg/conf"
	"profile_service/pkg/errors"
	"profile_service/pkg/models"
	"profile_service/pkg/providers"
)

type ResponseBody struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// Базовая ручка, чтобы ходить на auth_service/me
func ProfileDetailsHandler(config *conf.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		creds := &models.UserCredentials{}
		creds.AccessToken = r.Header.Get("Access")
		creds.RefreshToken = r.Header.Get("Refresh")

		authService := providers.NewHttpAuthServiceProvider(config)
		userData, err := getUserDataFromAuthService(r.Context(), creds, authService, *config)

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

		if resp, err := json.Marshal(userData); err != nil {
			makeInternalServerErrorResponse(&w)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(resp)
		}
	}
}

// Базовая ручка, чтобы ходить на auth_service/me
func ReceiversList(config *conf.Config, authService providers.AuthServiceProvider) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(ContextUserKey)
		body, _ := json.Marshal(user)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(body))
	}
	return IsAuthenticated(handler, config, authService)
}
