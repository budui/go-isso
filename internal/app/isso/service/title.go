package service

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
)

type TitleDiscriminator func(n *html.Node) (string, bool)

// TitleExtractor gets html title from remote page
type TitleExtractor struct {
	client http.Client
	tds    []TitleDiscriminator
}

func isTitleElement(n *html.Node) (title string, ok bool) {
	ok = n.Type == html.ElementNode && n.Data == "title"
	if ok {
		title = n.FirstChild.Data
	}
	return
}

func isDataTitle(n *html.Node) (title string, ok bool) {
	for _, a := range n.Attr {
		if a.Key == "data-title" {
			title = a.Val
			ok = true
			return
		}
	}
	return
}

// NewTitleExtractor makes extractor with cache. If memory cache failed, switching to no-cache
func NewTitleExtractor(client http.Client) *TitleExtractor {
	res := TitleExtractor{
		client: client,
		tds:    []TitleDiscriminator{isTitleElement, isDataTitle},
	}
	return &res
}

// Get page for url and return title
func (t *TitleExtractor) Get(url string) (string, error) {
	resp, err := t.client.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to load page %s (%v)", url, err)
	}
	defer resp.Body.Close() //nolint
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("can't load page %s, code %d", url, resp.StatusCode)
	}

	title, ok := t.getTitle(resp.Body)
	if !ok {
		return "", fmt.Errorf("can't get title for %s", url)
	}
	return title, nil
}

// get title from body reader, traverse recursively
func (t *TitleExtractor) getTitle(r io.Reader) (string, bool) {
	doc, err := html.Parse(r)
	if err != nil {
		log.Printf("[WARN] can't get header, %+v", err)
		return "", false
	}
	return t.traverse(doc)
}

func (t *TitleExtractor) traverse(n *html.Node) (string, bool) {
	for _, td := range t.tds {
		if t, ok := td(n); ok {
			return t, ok
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result, ok := t.traverse(c)
		if ok {
			return result, ok
		}
	}
	return "", false
}
