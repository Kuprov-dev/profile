package profile

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"profile_service/pkg/conf"
	"profile_service/pkg/db"
	requestErrors "profile_service/pkg/errors"
	"profile_service/pkg/models"
	"profile_service/pkg/providers"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Интерактор, который инкапсулирует логику работы с AuthServiceProvider
// и ходит за проверкой валидности токена в сервис auth
func checkUserIsAuthenticated(ctx context.Context, creds *models.UserCredentials, config *conf.Config, authService providers.AuthServiceProvider) (*models.UserAuthDetails, *models.RefreshedTokenCreds, error) {
	var err error
	var userAuthDetails *models.UserAuthDetails
	var refreshedTokens *models.RefreshedTokenCreds

	effector := func(ctx context.Context) error {
		userAuthDetails, refreshedTokens, err = authService.CheckUserIsAuthenticated(ctx, creds)
		log.Println("Effector ", *userAuthDetails, err)
		if err != nil {
			return err
		}
		return nil
	}

	effectorWithRetry := Retry(effector, config.AuthServiceRetries(), time.Duration(config.AuthServiceRetryDelay())*time.Millisecond)
	err = effectorWithRetry(ctx)

	return userAuthDetails, refreshedTokens, err
}

// Интерактор, который получает список рассылки юзера из UserDAO
func getUserReceivers(ctx context.Context, username string, userDAO db.UserDAO) (*models.UserRecievers, error) {
	userFromDB, err := userDAO.GetByUsername(ctx, username)

	if userFromDB == nil || err != nil {
		return nil, requestErrors.NewUserDAOError(requestErrors.UserNotFoundInDB, nil)
	}

	var receivers models.UserRecievers = models.UserRecievers{Receivers: userFromDB.Receivers}

	return &receivers, nil
}

// Интерактор который сопоставляет UUID из path и username, проверяет, что это это один и тот же юзер
// Служит для авторизации, кмк нужно переделать
func checkIsTheSameUser(ctx context.Context, userUUID primitive.ObjectID, username string, userDAO db.UserDAO) bool {
	user, err := userDAO.GetByUUID(ctx, userUUID)
	if user == nil || err != nil {
		return false
	}

	return user.Username == username
}

// Интерактор который добавляет айди юзера в список рассылки
func addReciever(ctx context.Context, userUUID primitive.ObjectID, receiverEmail string, userDAO db.UserDAO) error {
	user, err := userDAO.GetByUUID(ctx, userUUID)
	if err != nil {
		return err
	}

	if user == nil {
		return requestErrors.NewUserDAOError(requestErrors.UserNotFoundInDB, nil)
	}

	err = userDAO.AddReceiver(ctx, userUUID, receiverEmail)

	return err
}

// Интерактор для удаления юзера из списка рассылки
func removeReciever(ctx context.Context, userUUID primitive.ObjectID, receiverEmail string, userDAO db.UserDAO) error {
	user, err := userDAO.GetByUUID(ctx, userUUID)
	if err != nil {
		return err
	}
	if user == nil {
		return requestErrors.NewUserDAOError(requestErrors.UserNotFoundInDB, nil)
	}

	err = userDAO.RemoveReceiver(ctx, user.UUID, receiverEmail)

	return err
}

// Интерактор для сохранения шаблона и парсинга вунтренних параметров
func loadTemplateAndParseParams(ctx context.Context, templateData *models.HTMLTeplateCreateSchema, htmlTemplateDAO db.HTMLTemplateDAO) ([]string, error) {
	t := template.Must(template.New(templateData.Name).
		Parse(templateData.Template))
	params := ListTemplFields(t)

	_, err := htmlTemplateDAO.SaveTemplate(ctx, templateData, params, t)
	if err != nil {
		fmt.Println("???")
		return nil, err
	}
	return params, nil
}

func getUserDetails(ctx context.Context, username string, userDAO db.UserDAO) (*models.User, error) {
	var err error
	user, err := userDAO.GetByUsername(ctx, username)
	if user == nil {
		err = requestErrors.NewUserDAOError(requestErrors.UserNotFound, fmt.Errorf("User with username=%v not found.", username))
	}
	return user, err
}

func getTemplatesList(ctx context.Context, htmlTemplateDAO db.HTMLTemplateDAO) ([]*models.HTMLTeplate, error) {
	templates, err := htmlTemplateDAO.GetTemplatesList(ctx)

	if err != nil {
		return nil, err
	}

	return templates, nil
}

func getTemplateDetail(ctx context.Context, name string, htmlTemplateDAO db.HTMLTemplateDAO) (*models.HTMLTeplate, error) {
	templates, err := htmlTemplateDAO.GetTemplateByName(ctx, name)

	if err != nil {
		return nil, err
	}

	return templates, nil
}
