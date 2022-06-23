package collector

import (
	"fmt"
	"github.com/ZombieMInd/search-engine/internal/store/redisstore"
	"github.com/ZombieMInd/search-engine/internal/url_collector/domain"
	"github.com/go-redis/redis"
	"github.com/kelseyhightower/envconfig"
	"log"
)

func run() {
	cfg := &Config{}
	err := InitConfig(cfg)
	if err != nil {
		log.Fatal(err)
	}

	db, err := initDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	col := domain.NewCollector(db)

	for i := 0; i < 10; i++ {
		URL := db.PopDomain()
		domain.Collect(URL, db)
	}
}

func InitConfig(conf *Config) error {
	err := envconfig.Process("collector", conf)
	if err != nil {
		return fmt.Errorf("error while parsing env config: %s", err)
	}
	return nil
}

func initDB(cfg *Config) (*redisstore.Store, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return redisstore.New(client), nil
}
