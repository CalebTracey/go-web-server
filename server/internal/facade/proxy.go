package facade

import (
	"github.com/CalebTracey/go-web-server/server/config"
	"github.com/CalebTracey/go-web-server/server/internal/service/client"
	"github.com/CalebTracey/go-web-server/server/internal/service/environment"
	"github.com/CalebTracey/go-web-server/server/internal/service/server"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path/filepath"
)

type ProxyFacade interface {
	Client(urlPath string) (ClientResponse, error)
	Server() (*httputil.ReverseProxy, error)
	Environment() (environment.Response, error)
}

type Service struct {
	ClientSvc client.ServiceI
	ServerSvc server.ServiceI
	EnvSvc    environment.ServiceI
}

func NewService(appConfig *config.Config) (Service, error) {
	clientSvc, clientErr := client.InitializeClientSvc(appConfig)
	if clientErr != nil {
		return Service{}, clientErr
	}
	serverSvc, serverErr := server.InitializeServerSvc(appConfig)
	if clientErr != nil {
		return Service{}, serverErr
	}
	envSvc, envErr := environment.InitializeEnvService(appConfig)
	if clientErr != nil {
		return Service{}, envErr
	}
	return Service{
		ClientSvc: &clientSvc,
		ServerSvc: &serverSvc,
		EnvSvc:    &envSvc,
	}, nil
}

func (s Service) Client(urlPath string) (ClientResponse, error) {
	svc := s.ClientSvc.Client()
	var res ClientResponse

	path, err := filepath.Abs(urlPath)
	if err != nil {
		return res, err
	}
	res = ClientResponse{
		// prepend the path with the path to the static directory
		FilePath: filepath.Join(svc.StaticPath, path),
		// path to index.html for when file does not exist
		IndexPath:  filepath.Join(svc.StaticPath, svc.IndexPath),
		StaticPath: svc.StaticPath,
	}
	return res, nil
}

func (s Service) Server() (*httputil.ReverseProxy, error) {
	svc := s.ServerSvc.Server()

	target, svrErr := url.Parse(svc.TargetUrl)
	if svrErr != nil {
		return nil, svrErr
	}
	origin, uiErr := url.Parse(svc.HostUrl)
	if uiErr != nil {
		return nil, svrErr
	}
	director := func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", target.Host)
		req.Header.Add("X-Origin-Host", origin.Host)
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = target.Path

	}
	proxy := &httputil.ReverseProxy{Director: director}

	return proxy, nil
}

func (s Service) Environment() (environment.Response, error) {
	envResponse, envErr := s.EnvSvc.Set()
	if envErr != nil {
		return environment.Response{}, envErr
	}
	return envResponse, nil
}
