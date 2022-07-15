package collector

import (
	"fmt"
	"github.com/ZombieMInd/search-engine/internal/store/detabasestore"
	"github.com/ZombieMInd/search-engine/internal/store/redisstore"
	"github.com/ZombieMInd/search-engine/internal/url_collector/domain"
	"github.com/deta/deta-go/deta"
	"github.com/deta/deta-go/service/base"
	"github.com/go-redis/redis"
	"github.com/kelseyhightower/envconfig"
	"log"
)

func Run() {
	cfg := &Config{}
	err := InitConfig(cfg)
	if err != nil {
		log.Fatal(err)
	}

	db, err := initDetaBase(cfg)
	if err != nil {
		log.Fatal(err)
	}

	//col := domain.NewCollector(db)
	//
	//for i := 0; i < 10; i++ {
	//	URL := db.PopDomain()
	//	domain.Collect(URL, db)
	//}

	//db.AddDomain("somesite.com", false)
	//db.AddDomain("anothersite.com", false)
	//db.AddDomain("checked.com", true)
	//db.DeleteUnchecked("somesite.com")
	//
	//db.Add("some2", "data")

	d := db.GetDomainByKey("some2")
	fmt.Println(d)
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

func initDetaBase(conf *Config) (domain.Store, error) {
	d, err := deta.New(deta.WithProjectKey(conf.DetaBase.Key))
	if err != nil {
		return nil, err
	}

	// initialize with base name
	// returns ErrBadBaseName if base name is invalid
	db, err := base.New(d, "go-search-mind")
	if err != nil {
		return nil, err
	}

	store := detabasestore.New(db)

	return store, nil
}
