package main

import (
	"fmt"
	"net/url"
	"os"
	"path"
)

func qiwi(link string) (string, error) {
	parsedURL, err := url.Parse(link)
	if err != nil {
		return "Your file is invalid", err
	}

	// Check if the domain is "qiwi.gg"
	if parsedURL.Host == "qiwi.gg" {
		return "", fmt.Errorf("unsupported Link")
	}

	idUrl := path.Base(parsedURL.Path)
	endurl := fmt.Sprintf("https://spyderrock.com/%s", idUrl)
	return endurl, err
}
func main() {
	var url string
	fmt.Print("Enter URL: ")
	fmt.Scanln(&url)

	dlLink, err := qiwi(url)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Println("Download Link:", dlLink)
}
