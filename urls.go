package main

import (
	//"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getURLsFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	var urls []string
	doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		//fmt.Printf("href, exists = %s, %t\n", href, exists)
		if !exists {
			return
		}

		parsedURL, err := url.Parse(href)
		//fmt.Printf("parsedURL = %s\n", parsedURL)
		if err != nil {
			return
		}

		// Resolve relative URLs against the base URL
		absoluteURL := baseURL.ResolveReference(parsedURL)
		//fmt.Printf("absoluteURL = %s\n", absoluteURL)
		urls = append(urls, absoluteURL.String())
	})

	//fmt.Printf("urls = %s\n", urls)
	return urls, nil
}

func getImagesFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	var imageURLs []string
	doc.Find("img[src]").Each(func(_ int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if !exists {
			return
		}

		parsedURL, err := url.Parse(src)
		//fmt.Printf("parsedURL = %s\n", parsedURL)
		if err != nil {
			return
		}

		// Resolve relative URLs against the base URL
		absoluteURL := baseURL.ResolveReference(parsedURL)
		//fmt.Printf("absoluteURL = %s\n", absoluteURL)
		imageURLs = append(imageURLs, absoluteURL.String())
	})

	//fmt.Printf("imageURLs = %s\n", imageURLs)
	return imageURLs, nil
}
