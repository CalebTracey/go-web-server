package routes

import (
	"bytes"
	"encoding/json"
	"github.com/CalebTracey/go-web-server/server/internal/facade"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

type Handler struct {
	Service facade.ProxyFacade
}

type MyResponseWriter struct {
	http.ResponseWriter
	buf *bytes.Buffer
}

func (h Handler) InitializeRoutes() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	// Static routes
	r.Handle("/", h.ClientHandler())
	r.PathPrefix("/web/").Handler(http.StripPrefix("/web/", h.ClientHandler()))
	// Proxy route
	r.PathPrefix("/api/").Handler(h.ServiceHandler())
	// Health check
	r.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(map[string]bool{"ok": true})
		if err != nil {
			logrus.Errorln(err.Error())
			return
		}
	})
	return r
}

func (h Handler) ClientHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urlPath := r.URL.Path
		// build paths with request and config file
		res, svcErr := h.Service.Client(urlPath)
		if svcErr != nil {
			h.handleHttpError(w, svcErr, http.StatusBadRequest)
		}
		// check whether a file exists at the given path
		_, err := os.Stat(res.FilePath)
		if os.IsNotExist(err) {
			// file does not exist, serve index.html
			http.ServeFile(w, r, res.IndexPath)
			return
		} else if err != nil {
			h.handleHttpError(w, err, http.StatusInternalServerError)
			return
		}
		// default to using service.FileServer to serve the static dir
		http.FileServer(http.Dir(res.StaticPath)).ServeHTTP(w, r)
	}
}

func (h Handler) ServiceHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		serverProxy, svrErr := h.Service.Server()
		if svrErr != nil {
			logrus.Errorln(svrErr.Error())
		}
		serverProxy.Director(req)
		serverProxy.ServeHTTP(rw, req)
	}
}

func (h Handler) handleHttpError(w http.ResponseWriter, err error, status int) {
	logrus.Errorln(err.Error())
	http.Error(w, err.Error(), status)
}
