package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
)

type config struct {
	pages              map[string]PageData
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages           int
}

func main() {
	args := os.Args[1:]
	//fmt.Printf("args = %s\n", args)

	if len(args) < 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	} else if len(args) < 2 {
		fmt.Println("no maximum concurrency provided")
		os.Exit(1)
	} else if len(args) < 3 {
		fmt.Println("no maximum pages provided")
		os.Exit(1)
	}

	if len(args) > 3 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	rawBaseURL := args[0]
	fmt.Printf("starting crawl of: %s\n", rawBaseURL)

	parsedBaseURL, err := url.Parse(rawBaseURL)
	//fmt.Printf("base url = %s\n", rawBaseURL)
	if err != nil {
		fmt.Printf("error parsing base URL: %v\n", err)
		os.Exit(1)
	}

	maxConcurrency, err := strconv.Atoi(args[1])
	//fmt.Printf("maxConcurrency = %d\n", maxConcurrency)
	if err != nil {
		fmt.Printf("error parsing maxConcurrency: %v\n", err)
		os.Exit(1)
	}

	maxPagesArg, err := strconv.Atoi(args[2])
	//fmt.Printf("maxPages = %d\n", maxPagesArg)
	if err != nil {
		fmt.Printf("error parsing maxPages: %v\n", err)
		os.Exit(1)
	}

	// Initialize the configuration structure
	cfg := &config{
		pages:              make(map[string]PageData),
		baseURL:            parsedBaseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
		maxPages:           maxPagesArg,
	}

	// Kick off the first crawl
	cfg.wg.Add(1)
	go func() {
		defer cfg.wg.Done()
		cfg.crawlPage(rawBaseURL)
	}()

	// Wait for all concurrent goroutines to finish
	cfg.wg.Wait()

	// Print out the final results
	//fmt.Println("\n--- Final Crawl Results ---")
	//for normalizedURL, pageData := range cfg.pages {
	//	fmt.Printf("Found page: %s (Heading: %s) with %d outgoing links\n", normalizedURL, pageData.Heading, len(pageData.OutgoingLinks))
	//}

	writeJSONReport(cfg.pages, "report.json")
}

func (cfg *config) addPageVisit(normalizedURL string) (shouldVisit bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	// Check if the page is already in our map
	if _, visited := cfg.pages[normalizedURL]; visited {
		return false // already visited, skip it
	}

	// Mark it as visited (we can initialize it with an empty PageData for now)
	cfg.pages[normalizedURL] = PageData{}
	return true // first time seeing it, go ahead and crawl
}

func getHTML(rawURL string) (string, error) {
	// Create a new HTTP request
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set a User-Agent header to avoid being blocked by servers
	req.Header.Set("User-Agent", "BootCrawler/1.0")

	// Use an http.Client to fetch the webpage
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Return an error if the HTTP status code is 400+
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("bad status code: %d %s", resp.StatusCode, resp.Status)
	}

	// Return an error if the response content-type header is not text/html
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		return "", fmt.Errorf("invalid content-type: %s", contentType)
	}

	// Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Return the webpage's HTML as a string if successful
	return string(bodyBytes), nil
}

func (cfg *config) crawlPage(rawCurrentURL string) {
	// Return if max pages is exceeded
	cfg.mu.Lock()
	if len(cfg.pages) >= cfg.maxPages {
		cfg.mu.Unlock()
		return
	}
	cfg.mu.Unlock()

	// Parse and normalize the URL
	parsedURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("error parsing URL: %v\n", err)
		return
	}

	// Ensure we only crawl pages on the same domain/site
	if parsedURL.Hostname() != cfg.baseURL.Hostname() {
		return
	}

	normalizedURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("error normalizing URL: %v\n", err)
	}

	// Use our Mutex-protected helper method to check and claim the visit
	if !cfg.addPageVisit(normalizedURL) {
		return // if returns false, another goroutine already got to it
	}

	// If we got true, we are the lucky goroutine tasked with fetching it!
	fmt.Printf("crawling: %s\n", rawCurrentURL)

	htmlBody, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("error getting HTML for %s: %v\n", rawCurrentURL, err)
		return
	}

	// Extract links from the page
	pageData := extractPageData(htmlBody, rawCurrentURL)

	// Safely store the page data in the shared map using the mutex
	cfg.mu.Lock()
	cfg.pages[normalizedURL] = pageData
	cfg.mu.Unlock()

	// Loop through the outgoing links and spawn new concurrent goroutines
	for _, nextURL := range pageData.OutgoingLinks {
		cfg.wg.Add(1)
		go func(urlToCrawl string) {
			// Ensure token cleanup and WaitGroup decrement happen no matter what
			defer func() {
				<-cfg.concurrencyControl
				cfg.wg.Done()
			}()

			cfg.concurrencyControl <- struct{}{}

			cfg.crawlPage(urlToCrawl)
		}(nextURL)
	}
}
