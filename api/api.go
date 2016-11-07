package api

import (
	"encoding/json"
	// "errors"
	"fmt"
	// "io/ioutil"
	// "net"
	"net/http"
	// "net/url"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"

	"github.com/InVisionApp/kit-overwatch/config"
	"github.com/InVisionApp/kit-overwatch/deps"
)

const (
	DEFAULT_STATSD_RATE = 1.0
)

type Api struct {
	Config       *config.Config
	Version      string
	Dependencies *deps.Dependencies
}

type JSONStatus struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type detailedError struct {
	Message    string
	StatusCode int
}

type DetailedError struct {
	Error      error
	StatusCode int
}

type ExtendedJSONStatus struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

// Borrowed from http://laicos.com/writing-handsome-golang-middleware/
type Handler func(w http.ResponseWriter, r *http.Request) *DetailedError

func (a *Api) Handle(handlers map[string]Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for handlerName, handlerFunc := range handlers {
			var err *DetailedError

			// Record handler runtime
			func() {
				statusCode := "2xx"
				startTime := time.Now()

				if err = handlerFunc(w, r); err != nil {
					statusCode = strconv.Itoa(err.StatusCode)
					WriteJSONStatus(w, "error", err.Error.Error(), err.StatusCode)
				}

				// Record runtime metric
				go a.Dependencies.StatsD.TimingDuration(
					"handlers."+handlerName+".runtime",
					time.Since(startTime), // delta
					DEFAULT_STATSD_RATE,
				)

				// Record status code metric (default 2xx)
				go a.Dependencies.StatsD.Inc(
					"handlers."+handlerName+"."+statusCode,
					1,
					DEFAULT_STATSD_RATE,
				)
			}()

			// stop executing rest of the handlers if we encounter an error
			if err != nil {
				return
			}
		}
	})
}

func WriteJSONStatus(rw http.ResponseWriter, status, message string, statusCode int) {
	rw.Header().Set("Content-Type", "application/json")

	jsonData, _ := json.Marshal(&JSONStatus{
		Message: message,
		Status:  status,
	})

	rw.WriteHeader(statusCode)
	rw.Write(jsonData)
}

func New(cfg *config.Config, d *deps.Dependencies, version string) *Api {
	return &Api{
		Config:       cfg,
		Version:      version,
		Dependencies: d,
	}
}

func (a *Api) HomeHandler(rw http.ResponseWriter, r *http.Request) *DetailedError {
	fmt.Fprint(rw, "Refer to README.md for kit-overwatch API usage")
	return nil
}

func (a *Api) HealthHandler(rw http.ResponseWriter, r *http.Request) *DetailedError {
	rw.WriteHeader(200)
	rw.Write([]byte("Everything is peechy!"))
	return nil
}

func (a *Api) VersionHandler(rw http.ResponseWriter, r *http.Request) *DetailedError {
	fmt.Fprintf(rw, "auth-api: %v", a.Version)
	return nil
}

func (a *Api) Run() error {
	log.Infof("Starting API server on %v", a.Config.ListenAddress)

	routes := mux.NewRouter().StrictSlash(true)

	routes.Handle("/", a.Handle(map[string]Handler{
		"HomeHandler": a.HomeHandler,
	})).Methods("GET")

	routes.Handle("/version", a.Handle(map[string]Handler{
		"VersionHandler": a.VersionHandler,
	})).Methods("GET")

	routes.Handle("/healthcheck", a.Handle(map[string]Handler{
		"HealthHandler": a.HealthHandler,
	})).Methods("GET")

	return http.ListenAndServe(a.Config.ListenAddress, routes)
}
