package domain

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/ZombieMInd/search-engine/internal/url_collector/sitemap"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

func parsPage(doc *goquery.Document, baseURL *url.URL, pageCh chan *sitemap.Page, control chan int) {
	control <- 1

	defer func() {
		control <- -1
	}()

	page := &sitemap.Page{
		Url:      baseURL,
		Keywords: map[string]interface{}{},
	}

	doc.Find("title").Each(func(i int, s *goquery.Selection) {
		page.Title = s.Text()
	})

	doc.Find("meta[name='description']").Each(func(i int, s *goquery.Selection) {
		page.Description = s.Text()
	})

	doc.Find("h1").Each(func(i int, s *goquery.Selection) {
		words := strings.Split(s.Text(), " ")

		for _, w := range words {
			page.Keywords[w] = nil
		}
	})

	pageCh <- page
}

func collectURLS(doc *goquery.Document, baseURL *url.URL, externalURLsCh, internalURLsCh chan string, control chan int) {
	control <- 1

	defer func() {
		control <- -1
	}()

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		u, ok := s.Attr("href")
		if ok {
			if strings.HasPrefix(u, "/#") {
				return
			}

			parsed, errPars := url.Parse(u)
			if errPars != nil {
				return
			}

			if parsed.Host == "" || parsed.Host == baseURL.Host {
				internalURL := u
				if parsed.Host == "" {
					internalURL = baseURL.Scheme + "://" + baseURL.Host + u
				}
				internalURLsCh <- internalURL
			} else {
				externalURLsCh <- parsed.Scheme + "://" + parsed.Host
			}
		}
	})

}

func CollectFromPage(baseURL string, internalURLsCh, externalURLsCh chan string, pageCh chan *sitemap.Page, control chan int) {
	control <- 1
	ctx, cancel := context.WithTimeout(context.Background(), 2000*time.Millisecond)
	defer func() {
		control <- -1
		cancel()
	}()

	req, err := http.NewRequestWithContext(ctx, "GET", baseURL, nil)
	if err != nil {
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	defer func(Body io.ReadCloser) {
		errClose := Body.Close()
		if errClose != nil {
			logrus.Errorf("some closing err: %v", errClose)
		}
	}(res.Body)

	if res.StatusCode != 200 {
		return
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return
	}

	parsedBase, err := url.Parse(baseURL)
	if err != nil {
		return
	}

	go parsPage(doc, parsedBase, pageCh, control)

	doc.Find("meta[name='robots']").Each(func(i int, s *goquery.Selection) {
		content, ok := s.Attr("content")
		if ok {
			if content != "index, follow" {
				control <- -1
				return
			}
		}
	})

	go collectURLS(doc, parsedBase, externalURLsCh, internalURLsCh, control)

	return
}

type Store interface {
	SaveURLs(map[string]interface{}, *sync.WaitGroup)
	SaveSitemap(*sitemap.Sitemap, *sync.WaitGroup)
	DeleteUnchecked(d string)
	AddDomain(d string, checked bool)
}

func checkInternal(url string, urls map[string]interface{}) bool {
	l := len(urls)
	urls[url] = nil

	if l == len(urls) {
		return false
	}

	return true
}

func Collect(URL string, store Store) {
	internal := make(chan string, 5)
	external := make(chan string, 5)
	pages := make(chan *sitemap.Page, 5)
	control := make(chan int, 5)

	counter := 0

	sm := sitemap.New(URL)

	externalsMap := map[string]interface{}{}
	internalMap := map[string]interface{}{}

	logrus.Infof("collecting: %s", URL)

	go CollectFromPage(URL, internal, external, pages, control)

	store.DeleteUnchecked(URL)
	store.AddDomain(URL, true)

	for {
		select {
		case p := <-internal:
			if counter > 100 {
				time.Sleep(100 * time.Millisecond)
				internal <- p
				continue
			}

			if checkInternal(p, internalMap) {
				go CollectFromPage(p, internal, external, pages, control)
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
					go store.SaveURLs(externalsMap, wg)
				}

				wg.Add(1)
				go store.SaveSitemap(sm, wg)
				wg.Wait()

				return
			}
		}
	}

}
