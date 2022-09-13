package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

const maxWorkers = 5

func main() {
	urls := []string{"https://golang.org", "https://go.dev/tour", "https://go.dev/ref/spec", "https://golang.org", "https://go.dev/tour", "https://go.dev/ref/spec", "https://golang.org", "https://go.dev/tour", "https://go.dev/ref/spec", "https://golang.org", "https://go.dev/tour", "https://go.dev/ref/spec", "https://golang.org", "https://go.dev/tour", "https://go.dev/ref/spec", "https://golang.org", "https://go.dev/tour", "https://go.dev/ref/spec", "https://golang.org", "https://go.dev/tour", "https://go.dev/ref/spec", "https://golang.org", "https://go.dev/tour", "https://go.dev/ref/spec", "https://golang.org", "https://go.dev/tour", "https://go.dev/ref/spec", "https://golang.org", "https://go.dev/tour", "https://go.dev/ref/spec", "https://golang.org", "https://go.dev/tour", "https://go.dev/ref/spec", "https://golang.org", "https://go.dev/tour", "https://go.dev/ref/spec", "https://golang.org", "https://go.dev/tour", "https://go.dev/ref/spec", "https://golang.org", "https://go.dev/tour", "https://go.dev/ref/spec", "https://golang.org", "https://go.dev/tour", "https://go.dev/ref/spec", "https://golang.org", "https://go.dev/tour", "https://go.dev/ref/spec"}
	wg := sync.WaitGroup{}
	taskChan := make(chan string, maxWorkers)
	resultChan := make(chan int)

	httpClient := http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}
	var total int

	go func() {
		for i := 0; i < len(urls); i++ {
			total += <- resultChan
		}
	}()

	workersNum := len(urls)
	if workersNum > 5 {
		workersNum = maxWorkers
	}
	for i := 0; i < workersNum; i++ {
		wg.Add(1)
		go worker(&wg, taskChan, resultChan, &httpClient)
	}

	for _, url := range urls {
		taskChan <- url
	}
	close(taskChan)
	wg.Wait()
	close(resultChan)
	fmt.Println("Total: ", total)
}

func worker(wg *sync.WaitGroup, taskChan chan string, resultChan chan int, httpClient *http.Client) {
	defer wg.Done()

	for url := range taskChan {
		entriesCount := getEntriesCount(url, httpClient)
		fmt.Printf("%v: %v\n", url, entriesCount)
		resultChan <- entriesCount
	}
}

func getEntriesCount(u string, c *http.Client) int {
	resp, err := c.Get(u)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var entriesCount int
	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(body)
		entriesCount = strings.Count(bodyString, "Go")
	}
	return entriesCount
}