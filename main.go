package main

import (
	"log"
	"time"

	"github.com/praneethtkonda/go-web-crawler/worker"
)

func main() {
	// const ROOT_URL string = "http://crawler-test.com/"
	// const ROOT_URL string = "http://mockaroo.com/"
	// const ROOT_URL string = "https://webscraper.io/test-sites/e-commerce/static/product/126"
	const ROOT_URL string = "http://quotes.toscrape.com/"
	const NUM_WORKERS int = 100
	const FILEPATH string = "site_map_go.json"
	
	startTime := time.Now()
	worker.Start(ROOT_URL, NUM_WORKERS, FILEPATH)
	endTime := time.Since(startTime)

	log.Printf("Code Execution time: %v", endTime)
}