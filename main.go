package main

import (
	"flag"
	"fmt"
	"net/http"
	"regexp"

	"golang.org/x/net/html"
)

const baseURL = "https://en.wikipedia.org/wiki/"

func filterWikipediaURLs(urls []string) map[string]string {
	re := regexp.MustCompile(`^\/wiki\/`)
	ret := make(map[string]string)
	for _, url := range urls {
		if len(url) != 0 {
			found := re.FindStringIndex(url)
			if len(found) != 0 {
				ret[url[found[1]:]] = baseURL + url[found[1]:]
			}
		}
	}
	return ret
}

func gatherURLs(root *html.Node) []string {
	urls := make([]string, 100)
	rec(root, &urls)
	return urls
}

func rec(curr *html.Node, urls *[]string) {
	if curr.Type == html.ElementNode && curr.Data == "a" {
		for _, a := range curr.Attr {
			if a.Key == "href" {
				*urls = append(*urls, a.Val)
			}
		}
	}
	for c := curr.FirstChild; c != nil; c = c.NextSibling {
		rec(c, urls)
	}
}

func worker(url string, end string, depth int, maxDepth int, path []string, done chan<- []string) {
	if depth == maxDepth {
		return
	}
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return
	}
	urlsMap := filterWikipediaURLs(gatherURLs(doc))
	for childTitle, childURL := range urlsMap {
		if childTitle == end {
			done <- append(path, childTitle)
			return
		}
		go worker(childURL, end, depth+1, maxDepth, append(path, childTitle), done)
	}
}

func main() {
	fromArg := flag.String("from", "Go_(programming_language)", "Starting article. Sets to \"Go_(programming_language)\" by default")
	toArg := flag.String("to", "C_(programming_language)", "Destination article. Sets to \"C_(programming_language)\" by default")

	flag.Parse()

	fmt.Println(baseURL+*fromArg, " ", baseURL+*toArg)

	done := make(chan []string, 1)

	worker(baseURL+*fromArg, *toArg, 0, -1, []string{*fromArg}, done)
	fmt.Println(<-done)
}
