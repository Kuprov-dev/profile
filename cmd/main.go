package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"profile_service/pkg/conf"
	"profile_service/pkg/db"
	logging "profile_service/pkg/log"
	"profile_service/pkg/profile"
	"profile_service/pkg/providers"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func main() {
	// TODO dependencies
	// should be rewritten in get<Dependency>(...) style as in FastAPI
	config := conf.New()
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	logEntry := logrus.NewEntry(log)

	authService := providers.HttpAuthServiceProvider{Config: config}
	userDAO := db.InMemroyUserDAO{}
	htmlTemplateDAO := db.NewFSTemplateDAO("./html_templates/")

	r := mux.NewRouter()

	r.Handle("/i", profile.ProfileDetailsHandler(config, &authService))
	r.Handle("/receivers/{uuid}/", profile.ReceiversListHandler(config, &userDAO, &authService)).Methods("GET")
	r.Handle("/receivers/{uuid}/", profile.AddRecieverHandler(config, &userDAO, &authService)).Methods("POST")
	r.Handle("/receivers/{uuid}/", profile.RemoveRecieverHandler(config, &userDAO, &authService)).Methods("DELETE")
	r.Handle("/upload_template", profile.UploadHTMLTemplate(config, htmlTemplateDAO, &authService)).Methods("POST")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	handler := logging.LoggingMiddleware(logEntry)(r)
	s := &http.Server{
		Addr:    ":" + config.Port,
		Handler: handler,
	}
	defer s.Close()
	go func() {
		fmt.Println("Starting server")
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println(err)
			return
		}
	}()

	<-stop

	log.Fatal("Server shutdown")
}
