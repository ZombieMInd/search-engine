package main

import (
	"fmt"
	"github.com/ZombieMInd/search-engine/internal/store/redisstore"
	collector "github.com/ZombieMInd/search-engine/internal/url_collector"
	"github.com/go-redis/redis"
	"github.com/kelseyhightower/envconfig"
	"log"
	"time"
)

func main() {
	collection := map[string]int{}

	cfg := &collector.Config{}
	err := InitConfig(cfg)
	if err != nil {
		log.Fatal(err)
	}

	db, err := initDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	initialValue := db.GetFirstUnderValue(2)
	collection[initialValue] = 1

	if initialValue == "" {
		collection["https://moz.com/top500"] = 1
	}

	for {
		collection = collector.Collect(collection)
		log.Printf("collection size: %d\n", len(collection))

		saveCollection(collection, db)

		collection = map[string]int{
			db.GetFirstUnderValue(2): 1,
		}
		time.Sleep(5 * time.Second)
	}
	//
	//res := db.GetAll()
	//if res != nil {
	//	fmt.Println("done")
	//}
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
	AddDomain(string, int) error
}

func saveCollection(c map[string]int, db Store) {
	for name, val := range c {
		_ = db.AddDomain(name, val)
	}
}
