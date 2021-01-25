package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"mvdan.cc/xurls/v2"
)

const (
	currentWorker int8   = 15
	dirToSavePics string = "/home/eramirez/Downloads/pic/"
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

	var wg sync.WaitGroup
	wg.Add(len(filterLinks))
	linkChannel := make(chan string, currentWorker)

	go func() {
		for _, v := range filterLinks {
			linkChannel <- v
		}
	}()

	go func() {
		wg.Wait()
		close(linkChannel)
	}()

	for links := range linkChannel {

		go func(link string) {
			defer wg.Done()

			res, err := http.Get(link)
			CheckIfError(err)
			_, fileName := filepath.Split(link)

			defer res.Body.Close()
			file, err := os.Create(dirToSavePics + fileName)
			CheckIfError(err)
			defer file.Close()

			io.Copy(file, res.Body)
		}(links)
	}

}

func isAnImage(url string) bool {
	return strings.HasSuffix(url, ".png") || strings.HasSuffix(url, ".jpg")
}

//Function for filter an array based on a function
func filter(in []string, fn func(link string) bool) []string {
	result := make([]string, 0)
	linkSize := len(in)
	queueResult := make(chan string, currentWorker)
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
