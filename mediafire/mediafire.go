package main

import (
	"fmt"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/imroc/req/v3"
)

func mediafire(link string) (string, error) {
	client := req.C() // Create a new client

	resp, err := client.R().Get(link)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// for finding
	doc2, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	var href string
	doc2.Find("a.input.popsok[aria-label='Download file']").Each(func(i int, s *goquery.Selection) {
		href_exist, exist := s.Attr("href")
		if exist {
			href = href_exist
		} else {
			href = ""
		}
	})
	return href, nil
	// body, err := io.ReadAll(resp.Body) // read the response body
	// if err != nil {
	// 	return "", err // return an empty string and the error
	// }
	// file, err := os.Create("result.txt")
	// if err != nil {
	// 	return "", err
	// }
	// defer file.Close()
	// _, err = file.Write(body)
	// if err != nil {
	// 	return "", err
	// }
	// return "Successfully saved to result.txt", nil

}

func main() {
	var url string
	fmt.Print("Mediafire Link Generator\n")
	fmt.Print("Enter URL: ")
	fmt.Scanln(&url)

	dlLink, err := mediafire(url)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Println("Download Link:", dlLink)
}
