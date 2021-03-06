package parse

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func Parse(r io.Reader) ([]Link, error) {
	links := make([]Link, 0)
	parsedRootHTMLNode, htmlParseErr := html.Parse(r)
	if htmlParseErr != nil {
		return nil, htmlParseErr
	}
	links = iterateHTMLNodes(*parsedRootHTMLNode, links)
	return links, nil
}

func iterateHTMLNodes(htmlNode html.Node, links []Link) []Link {
	if htmlNode.Type == html.ElementNode && htmlNode.Data == "a" {
		link := Link{}
		for _, attribute := range htmlNode.Attr {
			if attribute.Key == "href" {
				link.Href = attribute.Val
			}
		}
		var textBuilder strings.Builder
		buildLinkText(&htmlNode, &textBuilder)
		link.Text = strings.TrimSpace(textBuilder.String())
		links = append(links, link)
	}
	for node := htmlNode.FirstChild; node != nil; node = node.NextSibling {
		links = iterateHTMLNodes(*node, links)
	}

	return links
}

func buildLinkText(node *html.Node, textBuilder *strings.Builder) {
	if node.Type == html.TextNode {
		textBuilder.WriteString(node.Data)
		if nextNode := node.FirstChild; nextNode != nil {
			buildLinkText(nextNode, textBuilder)
		} else if nextNode := node.NextSibling; nextNode != nil {
			buildLinkText(nextNode, textBuilder)
		}
		return
	} else if node.Type == html.ElementNode {
		if nextNode := node.FirstChild; nextNode != nil {
			buildLinkText(nextNode, textBuilder)
		} else if nextNode := node.NextSibling; nextNode != nil {
			buildLinkText(nextNode, textBuilder)
		}
		return
	}
}
