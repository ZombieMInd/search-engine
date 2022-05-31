package collector

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/ZombieMInd/search-engine/internal/constants"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func CollectFromPage(baseURL string, urlsCh chan string, done chan bool, level int) {
	if level <= 0 {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2000*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", baseURL, nil)
	if err != nil {
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		done <- true
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

	parsedBase, _ := url.Parse(baseURL)

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

			if parsed.Host == "" {
				u = parsedBase.Scheme + "://" + parsedBase.Host + u
			} else {
				urlsCh <- parsedBase.Scheme + "://" + parsed.Host
			}

			level--
			go CollectFromPage(u, urlsCh, done, level)
		}
	})

	return
}

func Collect(urls map[string]interface{}) map[string]interface{} {
	urlsCh := make(chan string, 1)
	done := make(chan bool)

	lastURL := constants.DefaultURLForCollection

	for i := range urls {
		lastURL = i
		delete(urls, i)
		break
	}

	go CollectFromPage(lastURL, urlsCh, done, 5)

	timeout := time.After(6 * time.Second)

	for {
		select {
		case val := <-urlsCh:
			urls[val] = nil
		case <-timeout:
			return urls
		case <-done:
			return urls
		}
	}
}
