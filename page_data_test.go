package main

import (
	"reflect"
	"testing"
)

func TestExtractPageData(t *testing.T) {
	inputURL := "https://crawler-test.com"
	inputBody := `<html><body>
        <h1>Test Title</h1>
        <p>This is the first paragraph.</p>
        <a href="/link1">Link 1</a>
        <img src="/image1.jpg" alt="Image 1">
    </body></html>`

	actual := extractPageData(inputBody, inputURL)

	expected := PageData{
		URL:            "https://crawler-test.com",
		Heading:        "Test Title",
		FirstParagraph: "This is the first paragraph.",
		OutgoingLinks:  []string{"https://crawler-test.com/link1"},
		ImageURLs:      []string{"https://crawler-test.com/image1.jpg"},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}
}
