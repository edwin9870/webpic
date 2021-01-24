package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"mvdan.cc/xurls/v2"
)

func main() {
	fmt.Printf("Hello world\n")
	url := "https://github.com/"

	response, err := http.Get(url)
	CheckIfError(err)

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	rxStrict := xurls.Strict()
	links := rxStrict.FindAllString(string(body), -1)

	filterLinks := filter(links, isAnImage)
	if filterLinks == nil {
	}

	log.Printf("Count of links found: %v", len(filterLinks))
}

func isAnImage(url string) bool {
	return strings.HasSuffix(url, ".png") || strings.HasSuffix(url, ".jpg")
}

//Function for filter an array based on a function
func filter(in []string, fn func(link string) bool) []string {
	result := make([]string, 0)
	for _, v := range in {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}
