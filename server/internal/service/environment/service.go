package environment

import (
	"github.com/CalebTracey/go-web-server/server/config"
	"github.com/sirupsen/logrus"
	"os"
)

//go:generate mockgen -source=service.go -destination=mockEnvService.go -package=environment
type ServiceI interface {
	Set() (Response, error)
}

type Service struct {
	environment string
	port        string
}

type Response struct {
	Environment string
	Port        string
}

func InitializeEnvService(appConfig *config.Config) (Service, error) {
	env := appConfig.Env
	port := appConfig.Port
	if env == "" {
		return Service{}, config.MissingField("environment")
	}
	if port == "" {
		return Service{}, config.MissingField("port")
	}
	return Service{
		environment: appConfig.Env,
		port:        appConfig.Port,
	}, nil
}

func (s *Service) Set() (Response, error) {
	var res Response
	port, portErr := s.setPort()
	if portErr != nil {
		return res, portErr
	}
	env, envErr := s.setEnvironment()
	if portErr != nil {
		return res, envErr
	}
	return Response{
		Environment: env,
		Port:        port,
	}, nil
}

func (s *Service) setPort() (string, error) {
	// sets environment port, try to get that first
	port := os.Getenv("PORT")
	if port == "" {
		port = s.port
		err := os.Setenv("PORT", port)
		if err != nil {
			return "", err
		}
	}

	return port, nil
}

func (s *Service) setEnvironment() (string, error) {
	envErr := os.Setenv("ENVIRONMENT", s.environment)
	if envErr != nil {
		logrus.Errorf("failed to set ENVIRONMENT; err: %v", envErr.Error())
		return "", envErr
	}
	return s.environment, nil
}
