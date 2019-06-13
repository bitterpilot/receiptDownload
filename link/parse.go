// Package link searches a html document for href's
// tutorial: Gophercises.com
package link

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

// Link is a simplified version of <a href="..."> link
type Link struct {
	Href string
	Text string
}

// Parse takes a HTML document and returns the links in a simplified struct.
func Parse(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	nodes := linkNodes(doc)
	var links []Link
	for _, node := range nodes {
		links = append(links, getLink(node))
	}
	return links, nil
}

func getLink(n *html.Node) Link {
	var link Link
	// Get Href
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			link.Href = attr.Val
		}
	}

	// Get text
	link.Text = getText(n)
	return link
}

func getText(n *html.Node) string {
	// check for other types of data, if it's just text ie: no child nodes
	// then immediately return
	if n.Type == html.TextNode {
		return n.Data
	}

	// Iterate through child nodes
	var text string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text += getText(c)
	}
	return strings.Join(strings.Fields(text), " ")
}

func linkNodes(n *html.Node) []*html.Node {
	if n.Type == html.ElementNode && n.Data == "a" {
		return []*html.Node{n}
	}
	var ret []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret = append(ret, linkNodes(c)...)
	}
	return ret
}
