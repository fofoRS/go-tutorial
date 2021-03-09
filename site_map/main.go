package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/fofoRS/go-tutorial/parse"
)

const (
	urlSet = "http://www.sitemaps.org/schemas/sitemap/0.9"
)

type SiteMap struct {
	XMLName xml.Name  `xml:"urlset"`
	Xmlns   string    `xml:"xmlns,attr"`
	Loc     []MapNode `xml:"url"`
}
type MapNode struct {
	Loc string `xml:"loc"`
}

func main() {
	targetURL := flag.String("site", "https://gophercises.com/", "Site you want to get the site map returned.")
	flag.Parse()
	siteMap := buildSiteMap(*targetURL)
	xmlEncoder := xml.NewEncoder(os.Stdout)
	if err := xmlEncoder.Encode(siteMap); err != nil {
		fmt.Println("Error")
	}
}

func buildSiteMap(targetURL string) SiteMap {
	htmlDocument, requestURL, err := getHTMLDocument(targetURL)
	if err != nil {
		panic(err)
	}
	siteMapNode := make(map[string]*MapNode)
	nodeSlice := make([]MapNode, 0)
	buildSiteNodes(htmlDocument, siteMapNode, *requestURL)
	for key := range siteMapNode {
		nodeSlice = append(nodeSlice, MapNode{key})
	}
	return SiteMap{Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9", Loc: nodeSlice}
}

func buildSiteNodes(htmlDocument string, siteMapNodes map[string]*MapNode, url url.URL) {
	links := getLinksFromHTML(htmlDocument, url)
	for _, link := range links {
		location := link.Href
		fmt.Printf("url: %s\n", location)
		if _, ok := siteMapNodes[location]; !ok {
			siteMapNodes[location] = &MapNode{Loc: location}
			document, requestURL, err := getHTMLDocument(location)
			if err != nil {
				continue
			}
			buildSiteNodes(document, siteMapNodes, *requestURL)
		}
	}
	return
}

func getLinksFromHTML(htmlDocument string, url url.URL) []parse.Link {
	htmlReader := strings.NewReader(htmlDocument)
	links, err := parse.Parse(htmlReader)
	if err != nil {
		panic(err)
	}
	normalizedLinks := make([]parse.Link, 0)
	for _, link := range links {
		normalizedLink(&link, url)
		if isLinkValid(link, url) {
			normalizedLinks = append(normalizedLinks, link)
		}
	}
	return normalizedLinks
}

func getHTMLDocument(url string) (string, *url.URL, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", nil, err
	}
	defer response.Body.Close()

	bodyAsBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", nil, err
	}
	return string(bodyAsBytes), response.Request.URL, nil
}

func normalizedLink(link *parse.Link, url url.URL) {
	switch {
	case strings.HasPrefix(link.Href, "http://") || strings.HasPrefix(link.Href, "https://"):
		return
	case strings.HasPrefix(link.Href, "/"):
		link.Href = url.Scheme + "://" + url.Host + link.Href
	default:
		link = &parse.Link{}
	}
}

func isLinkValid(link parse.Link, url url.URL) bool {
	targetURL, err := url.Parse(link.Href)
	if err != nil {
		return false
	}
	return (strings.HasPrefix(link.Href, "http://") || strings.HasPrefix(link.Href, "https://")) &&
		strings.Contains(targetURL.Host, url.Host)
}
