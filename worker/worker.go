package worker

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

// uRLCache is a thread-safe structure to store visited URLs.
type uRLCache struct {
	mu      sync.RWMutex
	visited map[string]bool
}

// Exists returns whether the URL exists in the cache.
func (u *uRLCache) Exists(key string) bool {
	u.mu.RLock()
	defer u.mu.RUnlock()
	return u.visited[key]
}

// Add adds the URL to the cache.
func (u *uRLCache) Add(key string) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.visited[key] = true
}

type siteMap struct {
	mu    sync.RWMutex
	pages map[string][]string
}

func (s *siteMap) Add(root, page string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pages[root] = append(s.pages[root], page)
}

func (s *siteMap) Export(filepath string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	file, err := os.Create(filepath)
	if err != nil {
		log.Fatalf("Error creating sitemap file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(s.pages)
}

func getLinks(resp *http.Response,
			  originURL *url.URL,
			  urlCache *uRLCache,
			  sitemap *siteMap) ([]string) {
	links := []string{}
	z := html.NewTokenizer(resp.Body)
    for {
        tokenType := z.Next()
        if tokenType == html.ErrorToken {
            return links
        }
        token := z.Token()
        if tokenType == html.StartTagToken && token.Data == "a" {
            for _, attr := range token.Attr {
                if attr.Key == "href" {
					attrURL := attr.Val

					parsedAttrURL, err := url.Parse(attrURL)
					if err == nil && parsedAttrURL.Host != "" && parsedAttrURL.Host != originURL.Host {
						continue
					}

					if parsedAttrURL == nil {
						continue
					}
					newURL := originURL.ResolveReference(parsedAttrURL).String()
					if !urlCache.Exists(newURL) {
						urlCache.Add(newURL)
						links = append(links, newURL)
					}
					sitemap.Add(originURL.String(), newURL)
                }
            }
        }
    }
}

func crawl(worker_seq int,
		   urlToBeParsed string,
		   client *http.Client,
		   urlCache *uRLCache,
		   sitemap *siteMap) ([]string, error){
	parsedURL, err := url.Parse(urlToBeParsed)
	if err != nil {
		log.Printf("[Worker %v] Error occurred while parsing URL: %v", worker_seq, err)
		return nil, err
	}
	log.Printf("Crawling URL: %v", parsedURL)

	resp, err := client.Get(urlToBeParsed)
	if err != nil {
		log.Printf("[Worker %v] Error occurred while fetching URL: %v", worker_seq, err)
		return nil, err
	}
	defer resp.Body.Close()
	return getLinks(resp, parsedURL, urlCache, sitemap), nil
}

func worker(seq int,
			bufferedQueueCh chan string,
			client *http.Client,
			urlCache *uRLCache,
			sitemap *siteMap,
			activeTasks *sync.WaitGroup) {
	// Worker code
	log.Printf("From worker %v", seq)
	for url := range bufferedQueueCh {
		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			activeTasks.Done()
			continue
		}
		log.Printf("[Worker %d] processing URL: %s", seq, url)
		links, err := crawl(seq, url, client, urlCache, sitemap)
		if err != nil {
			log.Printf("[Worker %d] Error from worker: %v", seq, err)
			activeTasks.Done()
			continue
		}

		for i := range links {
			activeTasks.Add(1)
			bufferedQueueCh <- links[i]
		}
		activeTasks.Done()
	}
}

func Start(root_url string, num_workers int, filepath string) {
	log.Println("Starting worker")
	var wg sync.WaitGroup
	var activeTasks sync.WaitGroup
	
	bufferedQueueCh := make(chan string, 1000)
	urlCache := uRLCache{visited: make(map[string]bool)}
	sitemap := siteMap{pages: make(map[string][]string)}

	// Creating the http client
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 30,
		},
	}

	for i := 0; i < num_workers; i++ {
		wg.Add(1)
		go func(seq int) {
			defer wg.Done()
			worker(seq, bufferedQueueCh, client, &urlCache, &sitemap, &activeTasks)
		}(i)
	}

	// Add the root URL to the queue and increment activeTasks
    activeTasks.Add(1)
    bufferedQueueCh <- root_url

	go func() {
        activeTasks.Wait()
        close(bufferedQueueCh)
    }()

	log.Println("Waiting for workers to finish")
	wg.Wait()
	log.Println("All workers finished")

	err := sitemap.Export(filepath)
	if err != nil {
		log.Fatalf("Error exporting sitemap: %v", err)
	} else {
		log.Printf("Sitemap exported to %v", filepath)
	}

	log.Println("Exiting worker")
}
