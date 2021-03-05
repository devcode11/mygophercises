package parser

import (
	"golang.org/x/net/html"
	"io"
	"fmt"
	"strings"
)

type Link struct {
	Href string
	Text string
}

// Links parses HTML from Reader r and returns a slice of Link
func Links(r io.Reader) ([]Link, error) {
	var links []Link = make([]Link, 0, 10)

	root, err := html.Parse(r)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	var searchAnchor func(*html.Node)
	var getText func(*html.Node) string

	getText = func(r *html.Node) string {
		text := strings.Builder{}
		for c:= r.FirstChild; c!= nil; c = c.NextSibling {
			if c.Type == html.TextNode {
				_, _ = text.WriteString(c.Data)
			} else if c.Type == html.ElementNode {
				text.WriteString(getText(c))
			}
		}
		return strings.TrimSpace(text.String())
	}

	searchAnchor = func(r *html.Node) {
		if r.Type == html.ElementNode && r.Data == "a" {
			link := Link{}
			for _, attr := range r.Attr {
				if attr.Key == "href" {
					link.Href = strings.TrimSpace(attr.Val)
					link.Text = strings.TrimSpace(getText(r))
					break
				}
			}
			links = append(links, link)
			return
		}

		for c := r.FirstChild; c != nil; c = c.NextSibling {
			searchAnchor(c)
		}
	}

	searchAnchor(root)
	return links, nil
}