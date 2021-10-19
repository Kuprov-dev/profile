package main

import (
	"fmt"
	"log"
	"net/http"
	"profile_service/pkg/conf"
	"profile_service/pkg/profile"
)

func main() {
	config := conf.New()

	http.HandleFunc("/i", profile.ProfileDetailsHandler(config))
	fmt.Println("Starting server")
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
