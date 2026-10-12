package main

import (
	//"fmt"
	"net/url"
)

type PageData struct {
	URL            string
	Heading        string
	FirstParagraph string
	OutgoingLinks  []string
	ImageURLs      []string
}

func extractPageData(htmlBody, pageURL string) PageData {
	parsedURL, err := url.Parse(pageURL)
	if err != nil {
		// Fallback or handle error if needed based on your implementation
		parsedURL = &url.URL{}
	}

	heading := getHeadingFromHTML(htmlBody)
	firstParagraph := getFirstParagraphFromHTML(htmlBody)

	outgoingLinks, err := getURLsFromHTML(htmlBody, parsedURL)
	if err != nil {
		outgoingLinks = []string{}
	}
	imageURLs, err := getImagesFromHTML(htmlBody, parsedURL)
	if err != nil {
		imageURLs = []string{}
	}

	return PageData{
		URL:            pageURL,
		Heading:        heading,
		FirstParagraph: firstParagraph,
		OutgoingLinks:  outgoingLinks,
		ImageURLs:      imageURLs,
	}
}
