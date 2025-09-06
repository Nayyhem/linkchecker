package checker

import (
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var client = &http.Client{
	Timeout: 5 * time.Second,
}

var rateLimiter = make(chan struct{}, 10)

func CheckLinkPage(link string) []string {
	resp, ok := linkTest(link)
	if !ok {
		return []string{"Lien inaccessible"}
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return []string{"Impossible de parser la page"}
	}

	links := findLinks(doc)
	return checkAllLinks(link, links)
}

func linkTest(link string) (*http.Response, bool) {
	resp, err := client.Get(link)
	if err != nil || resp.StatusCode != 200 {
		return nil, false
	}
	return resp, true
}

func isLinkAlive(link string) bool {
	if link == "" || link[0] == '#' ||
		strings.HasPrefix(link, "mailto:") ||
		strings.HasPrefix(link, "javascript:") ||
		strings.HasPrefix(link, "tel:") {
		return true
	}

	rateLimiter <- struct{}{}
	defer func() { <-rateLimiter }()

	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return false
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode < 400
}

func findLinks(doc *goquery.Document) []string {
	var links []string
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			links = append(links, href)
		}
	})
	return links
}

func makeAbsolute(base, href string) string {
	baseURL, err := url.Parse(base)
	if err != nil {
		return href
	}
	hrefURL, err := url.Parse(href)
	if err != nil {
		return href
	}
	return baseURL.ResolveReference(hrefURL).String()
}

func checkAllLinks(baseLink string, links []string) []string {
	var deadLinks []string
	deadChan := make(chan string)
	var wg sync.WaitGroup

	for _, l := range links {
		fullLink := makeAbsolute(baseLink, l)
		wg.Add(1)
		go func(fl string) {
			defer wg.Done()
			if !isLinkAlive(fl) {
				deadChan <- fl
			}
		}(fullLink)
	}

	go func() {
		wg.Wait()
		close(deadChan)
	}()

	for dl := range deadChan {
		deadLinks = append(deadLinks, dl)
	}

	return deadLinks
}
