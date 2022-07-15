package detabasestore

import (
	"fmt"
	"github.com/deta/deta-go/service/base"
	"log"
	"sync"
)

const (
	DomainsChecked   = "domains:checked"
	DomainsUnchecked = "domains:unchecked"
)

type Store struct {
	client *base.Base
}

type DomainData struct {
	Key  string `json:"key"`
	Data string
}

func New(c *base.Base) *Store {
	return &Store{client: c}
}

func (s *Store) SaveURLs(urls map[string]interface{}, wg *sync.WaitGroup) {
	defer wg.Done()

	for u := range urls {
		s.AddDomain(u, false)
	}
}

func (s *Store) DeleteUnchecked(d string) {
	err := s.client.Delete(DomainsUnchecked + ":" + d)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func (s *Store) Add(key, data string) {
	_, err := s.client.Put(&DomainData{
		Key:  key,
		Data: data,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Store) AddDomain(d string, checked bool) {
	keyPrefix := DomainsUnchecked
	if checked {
		keyPrefix = DomainsChecked
	}

	_, err := s.client.Put(&DomainData{
		Key:  keyPrefix + ":" + d,
		Data: d,
	})
	if err != nil {
		log.Fatal(err)
		return
	}
}

func (s *Store) GetDomainByKey(key string) string {
	var u DomainData

	// get item
	// returns ErrNotFound if no item was found
	err := s.client.Get(key, &u)
	if err != nil {
		log.Fatal(err)
	}

	return u.Data
}

func (s *Store) GetDomains() error {
	var results []*DomainData

	// fetch items
	_, err := s.client.Fetch(&base.FetchInput{
		Dest: &results,
	})
	if err != nil {
		return err
	}

	for i := range results {
		fmt.Println(*results[i])
	}

	return nil
}
