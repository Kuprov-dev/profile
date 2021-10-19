package profile

import (
	"net/http"
	"profile_service/pkg/conf"
)

func ProfileDetailsHandler(config *conf.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
