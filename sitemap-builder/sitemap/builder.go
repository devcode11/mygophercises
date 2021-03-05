package sitemap

import (
	"fmt"
	"io"
	"log"
	"net/http"
	// "net/url"
	"strings"

	"golang.org/x/net/html"
)

// Parses links for same same host
func getLinks(r io.Reader) ([]string, error) {

	root, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	links := make([]string, 0, 10)

	var searchAnchor func(*html.Node)

	searchAnchor = func(r *html.Node) {
		if r.Type == html.ElementNode && r.Data == "a" {
			for _, attr := range r.Attr {
				if attr.Key == "href" {
					links = append(links, strings.TrimSpace(attr.Val))
					break
				}
			}
			return
		}

		for c := r.FirstChild; c != nil; c = c.NextSibling {
			searchAnchor(c)
		}
	}

	searchAnchor(root)
	return links, nil
}

func Build(startUrl string) ([]string, error) {
	response, err := http.Get(startUrl)
	if err != nil {
		return nil, err
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("Build: Url %s returned %d", startUrl, response.StatusCode)
	}
	response.Body.Close()
	// scheme := response.Request.URL.Scheme
	baseUrl := response.Request.URL.Scheme + "://" + response.Request.URL.Hostname()
	// hostname := response.Request.URL.Hostname()

	if response.Request.URL.Hostname() == "" || response.Request.URL.Scheme == "" {
		return nil, fmt.Errorf("Cannot get hostname or scheme from url %s", startUrl)
	}
	startUrl = response.Request.URL.String()
	fmt.Println(response.Request.URL.RawPath)

	var siteMap = make(map[string]struct{})

	var crawl func(u string)

	valid := func(l string) string {
		if strings.HasPrefix(l, baseUrl) {
			return l
		} else if strings.HasPrefix(l, "/") {
			return baseUrl + l
		} else {
			return ""
		}
		
		/*
		purl, err := url.Parse(l)
		if err != nil {
			log.Println(err)
			return ""
		}

		if purl.Hostname() == hostname && purl.Scheme == scheme {
			return purl.String()
		} else if purl.Hostname() == "" && purl.Scheme == "" && purl.Path != "" {
			purl.Host = hostname
			purl.Scheme = scheme
			return purl.String()
		} else {
			return ""
		}
		*/
	}

	crawl = func(u string) {

		u = valid(u)
		if u == "" {
			return
		}

		if _, seen := siteMap[u]; seen {
			return
		}

		resp, err := http.Get(u)

		if err != nil {
			log.Println(err)
			return
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			log.Printf("crawl: Url %s returned %d", u, resp.StatusCode)
			return
		}

		siteMap[u] = struct{}{}
		// log.Println("Added", u)

		links, err := getLinks(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Println(err)
			return
		}

		for _, link := range links {
			crawl(link)
		}
	}

	crawl(startUrl)

	siteLinks := make([]string,0, len(siteMap))
	for k := range siteMap {
		siteLinks = append(siteLinks, k)
	}

	return siteLinks, nil
}
