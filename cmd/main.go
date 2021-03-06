package main

import (
	"context"
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
	db.Users = config.Database.Users

	fmt.Println(config.AuthService)

	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	logEntry := logrus.NewEntry(log)

	authService := providers.HttpAuthServiceProvider{Config: config}
	// userDAO := db.NewInMemoryUserDAO(config)
	// htmlTemplateDAO := db.NewFSTemplateDAO("./html_templates/")

	db.ConnectMongoDB(context.Background(), config)
	userDAO := db.NewMongoDBUserDAO(context.Background(), db.GetMongoDBConnection(), config)
	templateDAO := db.NewMongoDBTemplateDAO(context.Background(), db.GetMongoDBConnection(), config)

	r := mux.NewRouter()

	r.Handle("/i", profile.ProfileDetailsHandler(config, userDAO, &authService))
	r.Handle("/receivers/{uuid}/", profile.ReceiversListHandler(config, userDAO, &authService)).Methods(http.MethodGet)
	r.Handle("/receivers/{uuid}/", profile.AddRecieverHandler(config, userDAO, &authService)).Methods(http.MethodPost)
	r.Handle("/receivers/{uuid}/", profile.RemoveRecieverHandler(config, userDAO, &authService)).Methods(http.MethodDelete)
	r.Handle("/upload_template", profile.UploadHTMLTemplateHandler(config, templateDAO, &authService)).Methods(http.MethodPost)
	r.Handle("/templates", profile.HTMLTemplatesListHandler(config, templateDAO, &authService)).Methods(http.MethodGet)
	r.Handle("/templates/{name}/", profile.HTMLTemplateDetailHandler(config, templateDAO, &authService)).Methods(http.MethodGet)

	r.Handle("/test", profile.Test()).Methods(http.MethodGet)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	handler := logging.LoggingMiddleware(logEntry)(r)
	s := &http.Server{
		Addr:    config.ServerAddr(),
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
