package extract

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

// GetPageTitle do the GET request to url, then extract title
func GetPageTitle(ctx context.Context, host, uri string) (title string, newuri string, err error) {
	url := host + uri
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", "", fmt.Errorf("get title failed: do request failed %v, %w", url, err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return "", "", fmt.Errorf("get title failed: do request failed %v, %w", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", "", fmt.Errorf("get title failed: %s, %d", url, resp.StatusCode)
	}

	title, newuri, err = titleAndThreadURI(resp.Body, "Untitled", uri)
	if err != nil {
		return "", "", fmt.Errorf("get title failed: extract failed. %s, %w", url, err)
	}
	return
}

// titleAndThreadURI extract title and thread uri
// first try to get attribute ,
// if failed, get the <title> node content
func titleAndThreadURI(body io.Reader, defaultTitle string, defaultURI string) (title string, uri string, err error) {
	uri = defaultURI
	title = defaultTitle

	htmlRoot, err := html.Parse(body)
	if err != nil {
		return
	}

	issoRoot := getNodeByID(htmlRoot, "isso-thread")
	if issoRoot == nil {
		err = errors.New("can not find isso root in page")
		return
	}

	if u, ok := getAttrbyName(issoRoot, "data-isso-id"); ok {
		uri = u
	}
	title, ok := getAttrbyName(issoRoot, "data-title")
	if ok {
		return
	}

	title = defaultTitle
	if r := getNodeByTag(htmlRoot, "title"); r != nil {
		title = r.FirstChild.Data
		return
	}
	return
}

func getAttrbyName(n *html.Node, attrName string) (string, bool) {
	for _, a := range n.Attr {
		if a.Key == attrName {
			return a.Val, true
		}
	}
	return "", false
}

func getNodeByID(n *html.Node, id string) *html.Node {
	if n.Type == html.ElementNode && (n.Data == "div" || n.Data == "section") {
		for _, a := range n.Attr {
			if a.Key == "id" && id == a.Val {
				return n
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		tn := getNodeByID(c, id)
		if tn != nil {
			return tn
		}
	}
	return nil
}

func getNodeByTag(n *html.Node, tag string) *html.Node {
	if n.Type == html.ElementNode && strings.ToUpper(n.Data) == strings.ToUpper(tag) {
		return n
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		tn := getNodeByTag(c, tag)
		if tn != nil {
			return tn
		}
	}
	return nil
}
