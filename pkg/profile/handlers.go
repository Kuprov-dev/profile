package profile

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"profile_service/pkg/conf"
	"profile_service/pkg/db"
	"profile_service/pkg/errors"
	"profile_service/pkg/models"
	"profile_service/pkg/providers"
)

type ResponseBody struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// Базовая ручка, чтобы ходить на auth_service/me
func ProfileDetailsHandler(config *conf.Config, authService providers.AuthServiceProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		creds := &models.UserCredentials{}
		creds.AccessToken = r.Header.Get("Access")
		creds.RefreshToken = r.Header.Get("Refresh")

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

// Ручка списка рассылок для юзера. Возвращает айдишники юзеров
func ReceiversListHandler(config *conf.Config, userDAO db.UserDAO, authService providers.AuthServiceProvider) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		userValue := r.Context().Value(ContextUserDetailsKey)

		userDetails, ok := userValue.(*models.UserDetails)

		if !ok {
			makeInternalServerErrorResponse(&w)
			return
		}

		var err error
		var receivers *models.UserRecievers

		receivers, err = getUserReceivers(userDetails.Username, userDAO)

		if err != nil {
			ve, ok := err.(*errors.RequestError)

			if !ok {
				makeInternalServerErrorResponse(&w)
				return
			}

			switch {
			case ve.Errors&errors.UserNotFound != 0:
				makeNotFoundErrorResponse(&w)
				return
			default:
				makeBadRequestErrorResponse(&w, "bruuuh.")
				return
			}
		}

		if resp, err := json.Marshal(receivers); err != nil {
			makeInternalServerErrorResponse(&w)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(resp))
		}
	}

	// TODO подумать почему приходится инжектить несколько раз и исправить
	return IsAuthenticated(IsOwner(handler, config, userDAO, authService), config, userDAO, authService)
}

// Ручка добавления в список рассылки юзера
func AddRecieverHandler(config *conf.Config, userDAO db.UserDAO, authService providers.AuthServiceProvider) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		userValue := r.Context().Value(ContextUserIdKey)

		userId, ok := userValue.(int)
		if !ok {
			makeBadRequestErrorResponse(&w, "bruh")
			return
		}

		if !ok {
			makeBadRequestErrorResponse(&w, "bruh")
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			makeBadRequestErrorResponse(&w, "bruh")
		}

		var addReceiverData models.UserAddReceiver
		err = json.Unmarshal(body, &addReceiverData)

		if err != nil {
			makeBadRequestErrorResponse(&w, "bruh")
		}

		err = addReciever(userId, addReceiverData.ReceiverUsername, userDAO)

		if err != nil {
			makeBadRequestErrorResponse(&w, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}

	// TODO подумать почему приходится инжектить несколько раз и исправить
	return IsAuthenticated(IsOwner(handler, config, userDAO, authService), config, userDAO, authService)
}

// Ручка удаления из списока рассылки юзера
func RemoveRecieverHandler(config *conf.Config, userDAO db.UserDAO, authService providers.AuthServiceProvider) http.HandlerFunc {
	handler := func(w http.ResponseWriter, r *http.Request) {
		userValue := r.Context().Value(ContextUserIdKey)

		userId, ok := userValue.(int)
		if !ok {
			makeBadRequestErrorResponse(&w, "bruh")
			return
		}

		if !ok {
			makeBadRequestErrorResponse(&w, "bruh")
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			makeBadRequestErrorResponse(&w, "bruh")
		}

		var removeReceiverData models.UserRemoveReciever
		err = json.Unmarshal(body, &removeReceiverData)

		if err != nil {
			makeBadRequestErrorResponse(&w, "bruh")
		}

		err = removeReciever(userId, removeReceiverData.ReceiverUsername, userDAO)

		if err != nil {
			makeBadRequestErrorResponse(&w, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}

	// TODO подумать почему приходится инжектить несколько раз и исправить
	return IsAuthenticated(IsOwner(handler, config, userDAO, authService), config, userDAO, authService)
}
