package domain

import (
	"github.com/ZombieMInd/search-engine/internal/url_collector/sitemap"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Store interface {
	SaveURLs(map[string]interface{}, *sync.WaitGroup)
	DeleteUnchecked(d string)
	AddDomain(d string, checked bool)
	GetDomains() error
	GetDomainByKey(key string) string
	Add(key, data string)
}

type SearchStore interface {
	SaveSitemap(*sitemap.Sitemap, *sync.WaitGroup, chan error)
}

type Collector struct {
	store         Store
	searchStore   SearchStore
	tagsToAddToSM []string
	logger        logrus.Logger
}

func NewCollector(s Store, search SearchStore, tags []string) *Collector {
	return &Collector{
		store:         s,
		searchStore:   search,
		tagsToAddToSM: tags,
	}
}

func checkInternal(url string, urls map[string]interface{}) bool {
	l := len(urls)
	urls[url] = nil

	if l == len(urls) {
		return false
	}

	return true
}

func (c *Collector) Collect(URL string) {
	internal := make(chan string, 5)
	external := make(chan string, 5)
	pages := make(chan *sitemap.Page, 5)
	control := make(chan int, 5)
	errCh := make(chan error, 5)

	logger := c.logger.WithField("url", URL)

	proc := newDomainProcessor(pages, control, external, internal)

	counter := 0

	sm := sitemap.New(URL)

	externalsMap := map[string]interface{}{}
	internalMap := map[string]interface{}{}

	logrus.Infof("collecting: %s", URL)

	go proc.collectFromPage(URL)

	c.store.DeleteUnchecked(URL)
	c.store.AddDomain(URL, true)

	for {
		select {
		case p := <-internal:
			if counter > 100 {
				time.Sleep(100 * time.Millisecond)
				internal <- p
				continue
			}

			if checkInternal(p, internalMap) {
				go proc.collectFromPage(p)
			}

		case s := <-external:
			externalsMap[s] = nil

		case p := <-pages:
			sm.AddPage(*p)

		case num := <-control:
			counter += num
			logrus.Infof("goroutins counter: %d", counter)

			if counter <= 0 {
				wg := &sync.WaitGroup{}

				if externalsMap != nil {
					wg.Add(1)
					go c.store.SaveURLs(externalsMap, wg)
				}

				wg.Add(1)
				go c.searchStore.SaveSitemap(sm, wg, errCh)
				wg.Wait()

				select {
				case err := <-errCh:
					logger.Errorln(err)
				}

				return
			}
		}
	}

}
