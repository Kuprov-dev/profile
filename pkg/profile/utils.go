package profile

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"profile_service/pkg/errors"
	"profile_service/pkg/models"
	"text/template/parse"
	"time"
)

// Интерфейс оборачиваемой ф-ции для Retry
type Effoctor func(context.Context) error

// Реализация паттерна Retry
// Если ресурс недоступен, то пробуем подергать его еще
func Retry(effoctor Effoctor, retries int, delay time.Duration) Effoctor {
	return func(ctx context.Context) error {
		for r := 1; ; r++ {
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

// Вытаскивает все параметры из шаблона
func ListTemplFields(t *template.Template) []string {
	return listNodeFields(t.Tree.Root)
}

func listNodeFields(node parse.Node) []string {
	var res []string
	if node.Type() == parse.NodeAction {
		res = append(res, node.String())
	}

	if ln, ok := node.(*parse.ListNode); ok {
		for _, n := range ln.Nodes {
			res = append(res, listNodeFields(n)...)
		}
	}
	return res
}
