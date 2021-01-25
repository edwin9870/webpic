package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"

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
	linkSize := len(in)
	queueResult := make(chan string, 15)
	var wg sync.WaitGroup
	wg.Add(linkSize)
	for _, v := range in {

		go func(dataToValidate string, queueResult2 chan string, wg *sync.WaitGroup) {
			defer wg.Done()
			if fn(dataToValidate) {
				queueResult2 <- dataToValidate
			}
		}(v, queueResult, &wg)
	}

	go func() {
		wg.Wait()
		close(queueResult)
	}()
	for v := range queueResult {
		result = append(result, v)
	}
	log.Println("Finish filter")
	return result
}
