package client

import (
	"github.com/CalebTracey/go-web-server/server/config"
)

type ServiceI interface {
	Client() Response
}
type Service struct {
	staticPath string
	indexPath  string
}

type Response struct {
	StaticPath string
	IndexPath  string
}

func InitializeClientSvc(appConfig *config.Config) (Service, error) {
	staticPath := appConfig.ClientConfig.StaticPath
	indexPath := appConfig.ClientConfig.IndexPath
	if staticPath == "" || indexPath == "" {
		return Service{}, config.MissingField("paths")
	}
	return Service{
		staticPath: staticPath,
		indexPath:  indexPath,
	}, nil
}

func (s *Service) Client() Response {
	return Response{
		StaticPath: s.staticPath,
		IndexPath:  s.indexPath,
	}
}
