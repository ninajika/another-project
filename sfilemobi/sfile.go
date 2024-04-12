package main

import (
	"fmt"
	"html"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/imroc/req/v3"
)

// for typing
type Sfile struct {
	DownloadLink string
	FinalURL     string
	finalURL1    string
}

func sfile(url string) (string, string, error) {
	client := req.C() // Create a new client

	// First request to get the initial page
	resp, err := client.R().Get(url)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", "", err
	}

	var finalURL1, finalURL string
	doc.Find("#download").Each(func(i int, s *goquery.Selection) {
		finalURL1, _ = s.Attr("href")
	})

	// Second request to get the final page
	resp2, err := client.R().
		SetHeader("Referer", url).
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0").
		SetHeader("Cache-Control", "max-age=0").
		SetHeader("Upgrade-Insecure-Requests", "1").
		Get(finalURL1)
	if err != nil {
		return "", "", err
	}
	defer resp2.Body.Close()

	doc2, err := goquery.NewDocumentFromReader(resp2.Body)
	if err != nil {
		return "", "", err
	}

	html2, err := doc2.Html()
	if err != nil {
		return "", "", err
	}
	html2 = html.UnescapeString(html2)
	doc2.Find("#download").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}

		onclick, exists := s.Attr("onclick")
		if !exists {
			return
		}

		kValue := strings.TrimSuffix(strings.TrimPrefix(onclick, "location.href=this.href+'&k='+'"), ";return false;")
		finalURL = href + "&k=" + kValue
		// Remove the trailing single quote from the finalURL
		finalURL = strings.TrimSuffix(finalURL, "'")
	})
	return finalURL1, finalURL, nil
}

func main() {
	var url string
	fmt.Print("Enter URL: ")
	fmt.Scanln(&url)

	dlLink, finalURL, err := sfile(url)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Println("Download Link:", dlLink)
	fmt.Println("Final URL:", finalURL)
}
