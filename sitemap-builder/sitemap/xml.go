package sitemap

import (
	"encoding/xml"
	"io"
)

const urlsetXmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type lin struct {
	Val string `xml:"loc"`
}

type sitemapXml struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	Urls    []lin    `xml:"url"`
}

func ToXml(siteMap []string, w io.Writer) error {
	siteMapXml := sitemapXml{
		Urls:  make([]lin, len(siteMap)),
		Xmlns: urlsetXmlns,
	}
	for i, ll := range siteMap {
		siteMapXml.Urls[i] = lin{ll}
	}
	b, err := xml.MarshalIndent(siteMapXml, " ", "  ")
	if err != nil {
		return err
	}
	w.Write([]byte(xml.Header))
	w.Write(b)
	w.Write([]byte("\n"))
	return nil
}
