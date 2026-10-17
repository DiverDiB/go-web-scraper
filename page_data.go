package main

import (
	//bootdev run"fmt"
	"net/url"
)

type PageData struct {
	URL            string   `json:"url"`
	Heading        string   `json:"heading"`
	FirstParagraph string   `json:"first_paragraph"`
	OutgoingLinks  []string `json:"outgoing_links"`
	ImageURLs      []string `json:"image_urls"`
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
