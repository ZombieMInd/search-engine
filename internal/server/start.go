package server

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"

	"github.com/ZombieMInd/go-logger/internal/store"
	"github.com/ZombieMInd/go-logger/internal/store/redisstore"
	"github.com/ZombieMInd/go-logger/internal/store/sqlstore"
	"github.com/go-redis/redis"
	"github.com/gorilla/handlers"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
)

func Start(conf *Config) error {
	store, err := initStore(conf)
	if err != nil {
		log.Fatal(err)
	}

	srv := NewServer(store)
	srv.configLogger(conf)
	srv.InitServices(conf)
	initRouter(srv)

	return http.ListenAndServe(conf.BindAddr, srv)
}

func InitConfig(conf *Config) error {
	err := envconfig.Process("API", conf)
	if err != nil {
		return fmt.Errorf("error while parsing env config: %s", err)
	}
	return nil
}

func initRouter(s *server) {
	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
	s.router.HandleFunc("/log", s.handleLog()).Methods("POST")
	s.router.HandleFunc("/debug/pprof/", pprof.Index)
	s.router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	s.router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	s.router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	s.router.HandleFunc("/debug/pprof/trace", pprof.Trace)
}

func newDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func initStore(conf *Config) (store.Store, error) {
	if conf.StoreMode == "postgres" {
		db, err := newDB(conf.DBURL)
		if err != nil {
			return nil, err
		}
		return sqlstore.New(db), nil
	} else if conf.StoreMode == "redis" {
		client := redis.NewClient(&redis.Options{
			Addr:     conf.Redis.Host,
			Password: conf.Redis.Password,
			DB:       conf.Redis.DB,
		})

		_, err := client.Ping().Result()
		if err != nil {
			return nil, err
		}

		return redisstore.New(client), nil
	}
	return nil, errors.New("unknown store mode")
}
