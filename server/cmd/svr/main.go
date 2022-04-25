package main

import (
	"github.com/CalebTracey/go-web-server/server/config"
	"github.com/CalebTracey/go-web-server/server/internal/facade"
	"github.com/CalebTracey/go-web-server/server/internal/routes"
	"github.com/CalebTracey/go-web-server/server/internal/service"
	"github.com/NYTimes/gziphandler"
	"github.com/pkg/errors"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"os"
)

var configPath = "config.json"

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func main() {
	defer cheatDeath()

	appConfig := config.NewFromFile(configPath)
	appService, svcErr := facade.NewService(appConfig)
	if svcErr != nil {
		logrus.Panic(svcErr)
		panic(svcErr)
	}

	handler := routes.Handler{
		Service: appService,
	}

	env, envErr := appService.Environment()
	if envErr != nil {
		logrus.Errorf("environment error: %v", envErr.Error())
	}

	router := handler.InitializeRoutes()

	logrus.Infof("Current environment: %v", os.Getenv("ENVIRONMENT"))
	logrus.Fatal(service.ListenAndServe(env.Port, gziphandler.GzipHandler(cors.Default().Handler(router))))
}

func cheatDeath() {
	if r := recover(); r != nil {
		logrus.Errorf("very panic: %v,", r)
		if err, ok := r.(stackTracer); ok {
			logrus.Tracef("%+v", err.StackTrace()[0:2])
		}
	}
}
