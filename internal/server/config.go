package server

import "github.com/ZombieMInd/search-engine/internal/common"

type Config struct {
	Name     string `envconfig:"NAME" default:"search api"`
	BindAddr string `envconfig:"BIND_ADDR" default:"127.0.0.1:8080"`
	DB       common.DBConfig
	Redis    common.RedisConfig
}
