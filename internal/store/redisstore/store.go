package redisstore

import (
	"github.com/ZombieMInd/search-engine/internal/url_collector/sitemap"
	"github.com/go-redis/redis"
	"sync"
)

const (
	DomainsChecked   = "domains:checked"
	DomainsUnchecked = "domains:unchecked"
	Domains          = "domains"
	Sitemaps         = "sitemaps"
)

type Store struct {
	client *redis.Client
}

func New(c *redis.Client) *Store {
	return &Store{client: c}
}

func (s *Store) AddDomain(d string, checked bool) {
	if checked {
		s.client.SAdd(DomainsChecked, d)
	} else {
		s.client.SAdd(DomainsUnchecked, d)
	}
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

func (s *Store) SaveURLs(urls map[string]interface{}, wg *sync.WaitGroup) {
	defer wg.Done()

	for u := range urls {
		s.AddDomain(u, false)
	}
}

//easyjson:json
type RedisSitemap struct {
	Domain   string `json:"domain"`
	Keywords map[string][]RedisPage
}

//easyjson:json
type RedisPage struct {
	URL         string
	Title       string
	Description string
}

func (s *Store) SaveSitemap(sm *sitemap.Sitemap, wg *sync.WaitGroup) {
	defer wg.Done()

	toInsert := RedisSitemap{
		Domain:   sm.BaseURL.String(),
		Keywords: map[string][]RedisPage{},
	}

	for i := range sm.Keywords {
		var keywords []RedisPage
		for j := range sm.Keywords[i] {
			keywords = append(keywords, RedisPage{
				URL:         sm.Keywords[i][j].Url.String(),
				Title:       sm.Keywords[i][j].Title,
				Description: sm.Keywords[i][j].Description,
			})
		}
		toInsert.Keywords[i] = keywords
	}

	r, err := toInsert.MarshalJSON()
	if err != nil {
		return
	}

	s.client.Set("sitemaps:"+sm.BaseURL.String(), string(r), 0)
}

func (s *Store) DeleteUnchecked(d string) {
	s.client.SRem(DomainsUnchecked, d)
}
