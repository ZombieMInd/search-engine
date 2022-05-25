package server

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ZombieMInd/go-logger/internal/logger"
	"github.com/ZombieMInd/go-logger/internal/store"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	sessionName            = "gopherschool"
	ctxKeyRequestID ctxKey = iota
)

type ctxKey int8

type server struct {
	router   *mux.Router
	logger   *logrus.Logger
	services Services
	store    store.Store
}

type LogService interface {
	Save(*logger.LogRequest) error
}

type Services struct {
	Log LogService
}

func NewServer(store store.Store) *server {
	return &server{
		router: mux.NewRouter(),
		logger: logrus.New(),
		store:  store,
	}
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(ctxKeyRequestID),
		})
		logger.Infof("started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		var level logrus.Level
		switch {
		case rw.code >= 500:
			level = logrus.ErrorLevel
		case rw.code >= 400:
			level = logrus.WarnLevel
		default:
			level = logrus.InfoLevel
		}
		logger.Logf(
			level,
			"completed with %d %s in %v",
			rw.code,
			http.StatusText(rw.code),
			time.Since(start),
		)
	})
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func (s *server) configLogger(conf *Config) {
	formatter := &logrus.TextFormatter{
		FullTimestamp: true,
	}
	s.logger.SetFormatter(formatter)
}

func (s *server) InitServices(config *Config) error {
	repository := s.store.Log()
	s.services.Log = logger.NewService(repository)
	return nil
}
