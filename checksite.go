package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

const WORKERS = 25

func main() {
	site, _ := ioutil.ReadFile("sites.txt")
	siteFromByte := string(site)
	urls := strings.Split(siteFromByte, "\n")
	wg := new(sync.WaitGroup)
	in := make(chan string, 2*WORKERS)

	for i := 0; i < WORKERS; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for url := range in {
				checkSite(url)
			}
		}()
	}
	for _, url := range urls {
		if url != "" {
			in <- url
		}
	}
	close(in)
	wg.Wait()
}

func checkSite(url string) {
	now := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error getting site: %v", err)
	}
	defer resp.Body.Close()
	fmt.Println(time.Since(now), url)
}