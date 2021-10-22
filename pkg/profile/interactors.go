package profile

import (
	"context"
	"log"
	"profile_service/pkg/conf"
	"profile_service/pkg/db"
	requestErrors "profile_service/pkg/errors"
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

// Интерактор, который получает список рассылки юзера из UserDAO
func getUserReceivers(username string, userDAO db.UserDAO) (*models.UserRecievers, error) {
	userFromDB := userDAO.GetByUsername(username)

	if userFromDB == nil {
		return nil, requestErrors.NewUserDAOError(requestErrors.UserNotFoundInDB, nil)
	}

	var receivers models.UserRecievers = models.UserRecievers{Receivers: userFromDB.Receivers}

	return &receivers, nil
}

// Интерактор который сопоставляет Id из path и username, проверяет, что это это один и тот же юзер
// Служит для авторизации, кмк нужно переделать
func checkIsTheSameUser(userId int, username string, userDAO db.UserDAO) bool {
	user := userDAO.GetById(userId)
	if user == nil {
		return false
	}

	return user.Username == username
}

// Интерактор который добавляет айди юзера в список рассылки
func addReciever(userId int, receiverUsername string, userDAO db.UserDAO) error {
	user := userDAO.GetById(userId)
	if user == nil {
		return requestErrors.NewUserDAOError(requestErrors.UserNotFoundInDB, nil)
	}

	receiver := userDAO.GetByUsername(receiverUsername)
	if receiver == nil {
		return requestErrors.NewUserDAOError(requestErrors.ReceiverNotFoundInDB, nil)
	}

	if userId == receiver.ID {
		return requestErrors.NewUserDAOError(requestErrors.SameUser, nil)
	}

	err := userDAO.AddReceiver(user.ID, receiver.ID)

	return err
}

// Интерактор который добавляет айди юзера в список рассылки
func removeReciever(userId int, receiverUsername string, userDAO db.UserDAO) error {
	user := userDAO.GetById(userId)
	if user == nil {
		return requestErrors.NewUserDAOError(requestErrors.UserNotFoundInDB, nil)
	}

	receiver := userDAO.GetByUsername(receiverUsername)
	if receiver == nil {
		return requestErrors.NewUserDAOError(requestErrors.ReceiverNotFoundInDB, nil)
	}

	if userId == receiver.ID {
		return requestErrors.NewUserDAOError(requestErrors.SameUser, nil)
	}

	err := userDAO.RemoveReceiver(user.ID, receiver.ID)

	return err
}
