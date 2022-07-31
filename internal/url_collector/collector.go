package collector

import (
	"fmt"
	"github.com/ZombieMInd/search-engine/internal/store/redisstore"
	"github.com/ZombieMInd/search-engine/internal/store/searchstore"
	"github.com/ZombieMInd/search-engine/internal/url_collector/domain"
	"github.com/go-redis/redis"
	"github.com/kelseyhightower/envconfig"
	"github.com/meilisearch/meilisearch-go"
	"log"
)

func Run() {
	cfg := &Config{}
	err := InitConfig(cfg)
	if err != nil {
		log.Fatal(err)
	}

	db, err := initDBRedis(cfg)
	if err != nil {
		log.Fatal(err)
	}

	searchDB := initSearchDB(cfg)

	col := domain.NewCollector(db, searchDB, cfg.TagsToCollect)

	for i := 0; i < 10; i++ {
		URL := db.PopDomain()
		col.Collect(URL)
	}
}

func InitConfig(conf *Config) error {
	err := envconfig.Process("collector", conf)
	if err != nil {
		return fmt.Errorf("error while parsing env config: %s", err)
	}
	return nil
}

func initDBRedis(cfg *Config) (*redisstore.Store, error) {
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

func initSearchDB(cfg *Config) domain.SearchStore {
	client := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   cfg.SearchDB.Host,
		APIKey: cfg.SearchDB.APIKey,
	})

	db := searchstore.NewStore(client, cfg.SearchDB.SitemapIndex)

	return db
}
