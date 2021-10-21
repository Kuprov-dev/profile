package profile

import (
	"fmt"
	"net/http"
	"profile_service/pkg/conf"
)

// мидлварь чтобы проверить что юзер это самое
func IsAuthenticated(next http.HandlerFunc, config *conf.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Access: ", r.Header.Get("Access"))
		fmt.Println("Refresh: ", r.Header.Get("Refresh"))

		// creds := models.UserCredentials{
		// 	AccessToken:  r.Header.Get("Access"),
		// 	RefreshToken: r.Header.Get("Refresh"),
		// }

		next.ServeHTTP(w, r)
	})
}
