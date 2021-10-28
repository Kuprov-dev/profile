package profile

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/mail"
	"profile_service/pkg/conf"
	"profile_service/pkg/db"
	"profile_service/pkg/errors"
	"profile_service/pkg/models"
	"profile_service/pkg/providers"
	"strings"

	"github.com/google/uuid"
)

type ResponseBody struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// Базовая ручка, чтобы ходить на auth_service/me
func ProfileDetailsHandler(config *conf.Config, userDAO db.UserDAO, authService providers.AuthServiceProvider) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		userValue := r.Context().Value(ContextUserDetailsKey)

		userAuthDetails, ok := userValue.(*models.UserAuthDetails)

		if !ok {
			makeInternalServerErrorResponse(&w)
			return
		}

		userDetails, err := getUserDetails(r.Context(), userAuthDetails.Username, userDAO)

		if err != nil {
			makeBadRequestErrorResponse(&w, err.Error())
		}

		if resp, err := json.Marshal(userDetails); err != nil {
			makeInternalServerErrorResponse(&w)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(resp)
		}
	}

	isAuthenticatedMiddleware := IsAuthenticatedOrRefreshTokens(config, authService)
	return isAuthenticatedMiddleware(http.HandlerFunc(handler))
}

// Ручка списка рассылок для юзера. Возвращает айдишники юзеров
func ReceiversListHandler(config *conf.Config, userDAO db.UserDAO, authService providers.AuthServiceProvider) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		userValue := r.Context().Value(ContextUserDetailsKey)

		userAuthDetails, ok := userValue.(*models.UserAuthDetails)

		if !ok {
			makeInternalServerErrorResponse(&w)
			return
		}

		var err error
		var receivers *models.UserRecievers

		receivers, err = getUserReceivers(userAuthDetails.Username, userDAO)

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
	isAuthenticatedMiddleware := IsAuthenticatedOrRefreshTokens(config, authService)
	isOwnerMiddleware := IsOwner(config, userDAO, authService)

	return isAuthenticatedMiddleware(isOwnerMiddleware(http.HandlerFunc(handler)))
}

// Ручка добавления в список рассылки юзера
func AddRecieverHandler(config *conf.Config, userDAO db.UserDAO, authService providers.AuthServiceProvider) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		userValue := r.Context().Value(ContextUserUUIDKey)

		userUUID, ok := userValue.(uuid.UUID)
		if !ok {
			makeBadRequestErrorResponse(&w, "Get userUUID from context error.")
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			makeBadRequestErrorResponse(&w, "Read from body error.")
			return
		}

		var addReceiverData models.UserAddReceiver
		err = json.Unmarshal(body, &addReceiverData)

		if err != nil || addReceiverData.ReceiverEmail == "" {
			makeBadRequestErrorResponse(&w, "Receiver email is empty")
			return
		}

		if _, err := mail.ParseAddress(addReceiverData.ReceiverEmail); err != nil {
			makeBadRequestErrorResponse(&w, err.Error())
			return
		}

		err = addReciever(userUUID, strings.ToLower(addReceiverData.ReceiverEmail), userDAO)

		if err != nil {
			makeBadRequestErrorResponse(&w, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}

	// TODO подумать почему приходится инжектить несколько раз и исправить
	isAuthenticatedMiddleware := IsAuthenticatedOrRefreshTokens(config, authService)
	isOwnerMiddleware := IsOwner(config, userDAO, authService)

	return isAuthenticatedMiddleware(isOwnerMiddleware(http.HandlerFunc(handler)))
}

// Ручка удаления из списока рассылки юзера
func RemoveRecieverHandler(config *conf.Config, userDAO db.UserDAO, authService providers.AuthServiceProvider) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		userValue := r.Context().Value(ContextUserUUIDKey)

		userId, ok := userValue.(uuid.UUID)
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

		err = removeReciever(userId, strings.ToLower(removeReceiverData.ReceiverEmail), userDAO)

		if err != nil {
			makeBadRequestErrorResponse(&w, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}

	isAuthenticatedMiddleware := IsAuthenticatedOrRefreshTokens(config, authService)
	isOwnerMiddleware := IsOwner(config, userDAO, authService)

	return isAuthenticatedMiddleware(isOwnerMiddleware(http.HandlerFunc(handler)))
}

// Ручка удаления из списока рассылки юзера
func UploadHTMLTemplate(config *conf.Config, htmlTemplateDAO db.HTMLTemplateDAO, authService providers.AuthServiceProvider) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		var templateData models.HTMLTeplateCreateSchema
		if err := json.NewDecoder(r.Body).Decode(&templateData); err != nil {
			makeBadRequestErrorResponse(&w, "Error decoding data.")
			return
		}

		var params []string
		params, err := loadTemplateAndParseParams(r.Context(), &templateData, htmlTemplateDAO)
		if err != nil {
			makeBadRequestErrorResponse(&w, "Error decoding data.")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(models.HTMLTeplateParsedParamsResponse{Params: params})
		if err != nil {
			makeBadRequestErrorResponse(&w, "Encoding resposne error")
			return
		}
		w.WriteHeader(http.StatusOK)
	}

	isAuthenticatedMiddleware := IsAuthenticatedOrRefreshTokens(config, authService)

	return isAuthenticatedMiddleware(http.HandlerFunc(handler))
}
