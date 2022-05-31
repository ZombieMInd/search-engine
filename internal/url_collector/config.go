package collector

import "github.com/ZombieMInd/search-engine/internal/common"

type Config struct {
	Name  string `envconfig:"NAME" default:"collector"`
	DB    common.DBConfig
	Redis common.RedisConfig
}
