package collector

import "github.com/ZombieMInd/search-engine/internal/common"

type Config struct {
	Name          string `envconfig:"NAME" default:"collector"`
	TagsToCollect string `envconfig:"TAGS_TO_COLLECT" default:"h1,h2"`
	DB            common.DBConfig
	Redis         common.RedisConfig
	DetaBase      common.DetaBaseConfig
}
