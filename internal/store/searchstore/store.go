package searchstore

import (
	"fmt"
	"github.com/ZombieMInd/search-engine/internal/url_collector/sitemap"
	"github.com/meilisearch/meilisearch-go"
	"sync"
)

type Store struct {
	client       *meilisearch.Client
	sitemapIndex *meilisearch.Index
}

func NewStore(c *meilisearch.Client, sitemapIndex string) *Store {
	index := c.Index(sitemapIndex)
	return &Store{client: c, sitemapIndex: index}
}

func (s *Store) SaveSitemap(sm *sitemap.Sitemap, wg *sync.WaitGroup, errCh chan error) {
	defer wg.Done()
	res, err := s.sitemapIndex.AddDocuments(sm.Pages)
	if err != nil {
		errCh <- err
	}
	fmt.Println(res)
}
