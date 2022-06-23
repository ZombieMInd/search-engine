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
	"time"
)

type domainProcessor struct {
	pageCh                         chan *sitemap.Page
	control                        chan int
	externalURLsCh, internalURLsCh chan string
}

func newDomainProcessor(pageCh chan *sitemap.Page, control chan int, externalURLsCh, internalURLsCh chan string) *domainProcessor {
	return &domainProcessor{
		pageCh:         pageCh,
		control:        control,
		externalURLsCh: externalURLsCh,
		internalURLsCh: internalURLsCh,
	}
}

func (p *domainProcessor) parsPage(baseURL *url.URL, doc *goquery.Document) {
	p.control <- 1

	defer func() {
		p.control <- -1
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

	var text []string

	doc.Find("h1").Each(func(i int, s *goquery.Selection) {
		text = append(text, s.Text())
	})

	p.pageCh <- page
}

func (p *domainProcessor) collectURLS(baseURL *url.URL, doc *goquery.Document) {
	p.control <- 1

	defer func() {
		p.control <- -1
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
				p.internalURLsCh <- internalURL
			} else {
				p.externalURLsCh <- parsed.Scheme + "://" + parsed.Host
			}
		}
	})

}

func (p *domainProcessor) collectFromPage(baseURL string) {
	p.control <- 1
	ctx, cancel := context.WithTimeout(context.Background(), 2000*time.Millisecond)
	defer func() {
		p.control <- -1
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

	go p.parsPage(parsedBase, doc)

	doc.Find("meta[name='robots']").Each(func(i int, s *goquery.Selection) {
		content, ok := s.Attr("content")
		if ok {
			if content != "index, follow" {
				p.control <- -1
				return
			}
		}
	})

	go p.collectURLS(parsedBase, doc)

	return
}
