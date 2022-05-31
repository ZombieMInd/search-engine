package redisstore

import (
	"github.com/ZombieMInd/search-engine/internal/constants"
	"github.com/go-redis/redis"
	"strconv"
)

type Store struct {
	client *redis.Client
}

func New(c *redis.Client) *Store {
	return &Store{client: c}
}

func (s *Store) AddDomain(d string) error {
	s.client.SAdd("domains", d)

	return nil
}

func (s *Store) GetFirstUnderValue(limit int) string {
	keys, err := s.client.Keys("*").Result()
	if err != nil {
		return ""
	}

	for i := range keys {
		val := s.client.Get(keys[i]).Val()

		intVal, convErr := strconv.Atoi(val)
		if convErr != nil {
			continue
		}

		if intVal < limit {
			return keys[i]
		}
	}

	return constants.DefaultURLForCollection
}

func (s *Store) GetAll() []string {
	return s.client.SMembers("domains").Val()
}

func (s *Store) PopDomain() string {
	return s.client.SPop("domains").Val()
}

func (s *Store) GetRandomDomain() string {
	return s.client.SRandMember("domains").Val()
}

func (s *Store) GetDomainsCount() int64 {
	return s.client.SCard("domains").Val()
}
