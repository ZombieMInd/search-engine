package searchstore

import (
	"github.com/ZombieMInd/search-engine/internal/url_collector/sitemap"
	"github.com/meilisearch/meilisearch-go"
	"sync"
)

type Store struct {
	client       meilisearch.Client
	sitemapIndex string
}

func NewStore(c meilisearch.Client) *Store {
	return &Store{client: c}
}

func (s *Store) SaveSitemap(sm *sitemap.Sitemap, wg *sync.WaitGroup, errCh chan error) {
	defer wg.Done()

	index := s.client.Index(s.sitemapIndex)

	_, err := index.AddDocuments(sm.Pages)
	if err != nil {
		errCh <- err
	}
}
