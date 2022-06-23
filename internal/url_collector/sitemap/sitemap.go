package sitemap

import (
	"net/url"
)

type Page struct {
	Url         *url.URL
	Keywords    map[string]interface{}
	Title       string
	Description string
	Text        string
}

type Sitemap struct {
	BaseURL  *url.URL `json:"domain"`
	Keywords map[string][]Page
	Pages    []Page `json:"pages,omitempty"`
}

func New(URL string) *Sitemap {
	parsed, err := url.Parse(URL)
	if err != nil {
		return nil
	}

	return &Sitemap{
		BaseURL:  parsed,
		Keywords: map[string][]Page{},
		Pages:    []Page{},
	}
}

func (s *Sitemap) AddPage(p Page) {
	s.Pages = append(s.Pages, p)

	for w := range p.Keywords {
		if s.Keywords[w] == nil {
			s.Keywords[w] = []Page{}
		}

		s.Keywords[w] = append(s.Keywords[w], p)
	}
}
