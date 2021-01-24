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
	rePrefix := regexp.MustCompile(`^\/wiki\/`)
	reName := regexp.MustCompile(`^[^:]+$`)
	ret := make(map[string]string)
	for _, url := range urls {
		if len(url) != 0 {
			foundPrefix := rePrefix.FindStringIndex(url)
			if len(foundPrefix) != 0 {
				foundName := reName.FindString(url[foundPrefix[1]:])
				if foundName != "" {
					ret[foundName] = baseURL + foundName
				}
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

func worker(url string, end string, depth int, maxDepth int, path []string, result chan<- []string) {
	// fmt.Printf("%p\n", &path)
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
			result <- append(path, childTitle+" "+childURL)
			return
		}
		res := make([]string, len(path))
		copy(res, path)
		res = append(res, childTitle+" "+childURL)
		go worker(childURL, end, depth+1, maxDepth, res, result)
	}
}

func main() {
	fromArg := flag.String("from", "Go_(programming_language)", "Starting article. Sets to \"Go_(programming_language)\" by default")
	toArg := flag.String("to", "C_(programming_language)", "Destination article. Sets to \"C_(programming_language)\" by default")

	flag.Parse()

	fmt.Println(baseURL+*fromArg, " ", baseURL+*toArg)

	result := make(chan []string, 1)
	worker(baseURL+*fromArg, *toArg, 0, -1, []string{*fromArg}, result)
	fmt.Println(<-result)
}
