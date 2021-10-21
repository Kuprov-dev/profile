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

func makeInternalServerErrorResponse(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "application/json")
	(*w).WriteHeader(http.StatusInternalServerError)
	body, _ := json.Marshal(ResponseBody{Status: 500, Message: "Internal error."})
	(*w).Write(body)
}

func makeBadGatewayErrorResponse(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "application/json")
	(*w).WriteHeader(http.StatusBadGateway)
	body, _ := json.Marshal(ResponseBody{Status: 502, Message: "Bad gateway."})
	(*w).Write(body)
}

func makeServiceUnavailableErrorResponse(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "application/json")
	(*w).WriteHeader(http.StatusServiceUnavailable)
	body, _ := json.Marshal(ResponseBody{Status: 503, Message: "External service is unavailable."})
	(*w).Write(body)
}

func makeBadRequestErrorResponse(w *http.ResponseWriter, errMsg string) {
	(*w).Header().Set("Content-Type", "application/json")
	(*w).WriteHeader(http.StatusBadRequest)
	body, _ := json.Marshal(ResponseBody{Status: 400, Message: errMsg})
	(*w).Write(body)
}

func makeUnathorisedErrorResponse(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "application/json")
	(*w).WriteHeader(http.StatusUnauthorized)
	body, _ := json.Marshal(ResponseBody{Status: 401, Message: "Unathorised error."})
	(*w).Write(body)
}

func makeForbiddenErrorResponse(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "application/json")
	(*w).WriteHeader(http.StatusUnauthorized)
	body, _ := json.Marshal(ResponseBody{Status: 403, Message: "Forbidden error."})
	(*w).Write(body)
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
				makeBadRequestErrorResponse(&w, "Something goes wrong.")
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
