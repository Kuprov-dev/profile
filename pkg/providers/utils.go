package providers

import (
	"net/http"
	"profile_service/pkg/models"
)

func getTokenCookiesFromResponse(w *http.Response) *models.RefreshedTokenCreds {
	var refreshedTokens models.RefreshedTokenCreds
	for _, cookie := range (*w).Cookies() {
		switch {
		case cookie.Name == "Access" && cookie.Value != "":
			refreshedTokens.AccessToken = cookie.Value
			refreshedTokens.AccessExpirationTime = cookie.Expires
		case cookie.Name == "Refresh" && cookie.Value != "":
			refreshedTokens.RefreshedToken = cookie.Value
			refreshedTokens.RefreshExpirationTime = cookie.Expires
		}
	}
	if refreshedTokens.AccessToken != "" && refreshedTokens.RefreshedToken != "" {
		return &refreshedTokens
	}
	return nil
}
