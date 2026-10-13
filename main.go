package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	args := os.Args[1:]
	//fmt.Printf("args = %s\n", args)

	if len(args) < 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	}

	if len(args) > 1 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	baseURL := args[0]
	fmt.Printf("starting crawl of: %s\n", baseURL)

	pages := make(map[string]int)
	crawlPage(baseURL, baseURL, pages)

	fmt.Println("\n--- Final Crawl Results ---")
	for page, count := range pages {
		fmt.Printf("Found %d internal links to %s\n", count, page)
	}
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

func crawlPage(rawBaseURL string, rawCurrentURL string, pages map[string]int) {
	// Make sure rawCurrentURL is on the same doman as rawBaseURL
	baseURLParsed, err := url.Parse(rawBaseURL)
	//fmt.Printf("baseURLParsed = %s\n", baseURLParsed)
	if err != nil {
		fmt.Printf("error parsing base URL: %v\n", err)
		return
	}

	currentURLParsed, err := url.Parse(rawCurrentURL)
	//fmt.Printf("currentURLParsed = %s\n", currentURLParsed)
	if err != nil {
		fmt.Printf("error parsing current URL: %v\n", err)
		return
	}

	if baseURLParsed.Hostname() != currentURLParsed.Hostname() {
		fmt.Println("URLs are not on the same domain")
		return
	}

	// Get a normalized version of the rawCurrentURL
	normalizedCurrent, err := normalizeURL(rawCurrentURL)
	//fmt.Printf("normalizedCurrent = %s\n", normalizedCurrent)
	if err != nil {
		fmt.Printf("error normalizing URL '%s': %v\n", rawCurrentURL, err)
		return
	}

	// If already visited, increment count and return
	if _, exists := pages[normalizedCurrent]; exists {
		pages[normalizedCurrent]++
		return
	}

	// Otherwise, add to map with a count of 1
	pages[normalizedCurrent] = 1

	// Get the HTML from the current URL with a live print statement
	fmt.Printf("crawling: %s\n", rawCurrentURL)
	htmlBody, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("error getting HTML for %s: %v\n", rawCurrentURL, err)
		return
	}

	// Parse the rawBaseURL into a *url.URL object
	baseURLParsedObj, err := url.Parse(rawBaseURL)
	if err != nil {
		fmt.Printf("error parsing base URL for link extraction: %v\n", err)
		return
	}

	// Get all URLs from the response body HTML using the parsed URL
	nextURLs, err := getURLsFromHTML(htmlBody, baseURLParsedObj)
	//fmt.Printf("nextURLs = %v\n", nextURLs)
	if err != nil {
		fmt.Printf("error extracting URLs from %s: %v\n", rawCurrentURL, err)
		return
	}

	// Recursively crawl each URL found on the page
	for _, nextURL := range nextURLs {
		crawlPage(rawBaseURL, nextURL, pages)
	}
}
