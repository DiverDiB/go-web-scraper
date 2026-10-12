package main

import (
	//"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getHeadingFromHTML(html string) string {
	//fmt.Print("getHeadingFromHTML\n")
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}

	// Try h1 first
	heading := doc.Find("h1").First().Text()
	//fmt.Printf("heading = %s\n", heading)
	if heading != "" {
		return strings.TrimSpace(heading)
	}

	// Fallback to h2
	heading = doc.Find("h2").First().Text()
	//fmt.Printf("heading = %s\n", heading)
	if heading != "" {
		return strings.TrimSpace(heading)
	}

	return ""
}
func getFirstParagraphFromHTML(html string) string {
	//fmt.Print("getFirstParagraphFromHTML\n")
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}
	// Check if a <main> tag exists, and look for a <p> inside it first
	mainSelection := doc.Find("main")
	//fmt.Printf("mainSelection = %s\n", mainSelection.Text())
	if mainSelection.Length() > 0 {
		p := mainSelection.Find("p").First().Text()
		if p != "" {
			return strings.TrimSpace(p)
		}
	}
	// Fall back to the very first <p> in the document
	p := doc.Find("p").First().Text()
	//fmt.Printf("p = %s\n", p)
	if p != "" {
		return strings.TrimSpace(p)
	}
	return ""
}
