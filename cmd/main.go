package main

import (
	"fmt"
	"log"
	"net/http"
	"profile_service/pkg/conf"
	"profile_service/pkg/profile"

	"github.com/gorilla/mux"
)

func main() {
	config := conf.New()

	r := mux.NewRouter()

	r.HandleFunc("/i", profile.ProfileDetailsHandler(config))
	r.HandleFunc("/receivers/{id:[0-9]+}", profile.ReceiversList(config)).Methods("GET")

	fmt.Println("Starting server")
	log.Fatal(http.ListenAndServe("localhost:8000", r))
}
