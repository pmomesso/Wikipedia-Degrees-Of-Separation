package main

import (
	"flag"
	"fmt"
)

const baseURL = "https://en.wikipedia.org/wiki/"

func main() {
	fromArg := flag.String("from", "Go_(programming_language)", "Starting article. Sets to \"Go_(programming_language)\" by default")
	toArg := flag.String("to", "C_(programming_language)", "Destination article. Sets to \"C_(programming_language)\" by default")

	flag.Parse()

	fmt.Println(baseURL+*fromArg, " ", baseURL+*toArg)

}
