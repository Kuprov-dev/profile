package main

import (
	"fmt"
	"log"
	"net/http"
	"profile_service/pkg/conf"
	"profile_service/pkg/db"
	"profile_service/pkg/profile"
	"profile_service/pkg/providers"

	"github.com/gorilla/mux"
)

func main() {
	// TODO dependencies
	// should be rewritten in get<Dependency>(...) style as in FastAPI
	config := conf.New()
	authService := providers.HttpAuthServiceProvider{Config: config}
	userDAO := db.InMemroyUserDAO{}

	r := mux.NewRouter()

	r.HandleFunc("/i", profile.ProfileDetailsHandler(config))
	r.HandleFunc("/receivers/{id:[0-9]+}", profile.ReceiversList(config, &userDAO, &authService)).Methods("GET")

	fmt.Println("Starting server")
	log.Fatal(http.ListenAndServe("localhost:8000", r))
}
