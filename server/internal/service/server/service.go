package server

import "github.com/CalebTracey/go-web-server/server/config"

type ServiceI interface {
	Server() Response
}
type Service struct {
	hostUrl   string
	targetUrl string
}

type Response struct {
	HostUrl   string
	TargetUrl string
}

func InitializeServerSvc(appConfig *config.Config) (Service, error) {
	hostUrl := appConfig.ClientConfig.Url
	targetUrl := appConfig.ClientConfig.Url
	if hostUrl == "" || targetUrl == "" {
		return Service{}, config.MissingField("urls")
	}
	return Service{
		hostUrl:   appConfig.ClientConfig.Url,
		targetUrl: appConfig.ClientConfig.Url,
	}, nil
}

func (s *Service) Server() Response {
	return Response{
		HostUrl:   s.hostUrl,
		TargetUrl: s.targetUrl,
	}
}
