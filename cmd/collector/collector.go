package main

import (
	"fmt"
	"github.com/ZombieMInd/search-engine/internal/constants"
	"github.com/ZombieMInd/search-engine/internal/store/redisstore"
	collector "github.com/ZombieMInd/search-engine/internal/url_collector"
	"github.com/go-redis/redis"
	"github.com/kelseyhightower/envconfig"
	"log"
	"time"
)

func main() {

	for i := 0; i < 50; i++ {
		go run()
	}

	go func() {
		cfg := &collector.Config{}
		err := InitConfig(cfg)
		if err != nil {
			log.Fatal(err)
		}

		db, err := initDB(cfg)
		if err != nil {
			log.Fatal(err)
		}
		for {
			log.Printf("total count: %d", db.GetDomainsCount())
			time.Sleep(5 * time.Second)
		}
	}()

	run()
}

func run() {
	collection := map[string]interface{}{}

	cfg := &collector.Config{}
	err := InitConfig(cfg)
	if err != nil {
		log.Fatal(err)
	}

	db, err := initDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	collection = initCollection(db)

	for {
		collection = collector.Collect(collection)
		//log.Printf("collected %d domains\n", len(collection))

		if len(collection) > 0 {
			saveCollection(collection, db)
		}

		collection = initCollection(db)
	}
}

func initCollection(db Store) map[string]interface{} {
	initialValue := db.GetRandomDomain()
	if initialValue == "" {
		initialValue = constants.DefaultURLForCollection
	}

	//log.Printf("collection initialized with %s\n", initialValue)

	return map[string]interface{}{
		initialValue: nil,
	}
}

func InitConfig(conf *collector.Config) error {
	err := envconfig.Process("collector", conf)
	if err != nil {
		return fmt.Errorf("error while parsing env config: %s", err)
	}
	return nil
}

func initDB(cfg *collector.Config) (*redisstore.Store, error) {
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

type Store interface {
	AddDomain(string) error
	GetRandomDomain() string
}

func saveCollection(c map[string]interface{}, db Store) {
	for name := range c {
		_ = db.AddDomain(name)
	}
}
