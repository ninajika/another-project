package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"regexp"

	"github.com/imroc/req/v3"
)

func extractWebToken(input string, key string) (string, error) {
	re := regexp.MustCompile(key + `: "([^"]+)"`)
	match := re.FindStringSubmatch(input)
	if len(match) < 2 {
		return "", fmt.Errorf("key %s not found", key)
	}

	return match[1], nil
}

func gofile(link string) (string, string, error) {
	client := req.C() // Create a new client

	// make a guest account
	resp, err := client.R().Post("https://api.gofile.io/accounts")
	if err != nil {
		return "guest account changed again", "", err
	}

	defer resp.Body.Close()

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "guest account changed again", "", err
	}

	// just get the "token" value
	data, ok := response["data"].(map[string]interface{})
	if !ok {
		return "Data not found in response", "", errors.New("data not found in guest account, maybe something goes wrong")
	}

	tokenValue, tokenExists := data["token"].(string)
	if !tokenExists {
		return "Token not found in response", "", errors.New("token not found, maybe something goes wrong")
	}

	// for getting id url
	parsedURL, err := url.Parse(link)
	if err != nil {
		return "Your file is invalid", "", err
	}

	idUrl := path.Base(parsedURL.Path)

	// for getting websiteToken
	resp2, err := client.R().Get("https://gofile.io/dist/js/alljs.js")
	if err != nil {
		return "WebsiteToken Getter Broken", "", err
	}

	defer resp2.Body.Close()

	wtGet, err := io.ReadAll(resp2.Body)
	if err != nil {
		return "", "", err
	}

	wtValue, err := extractWebToken(string(wtGet), "wt")
	if err != nil {
		return "", "", err
	}

	// get the file
	url_dl := fmt.Sprintf("https://api.gofile.io/contents/%s?wt=%s", idUrl, wtValue)
	token_dl := fmt.Sprintf("Bearer %s", tokenValue)

	resp3, err := client.R().
		SetHeader("Authorization", token_dl).
		SetHeader("origin", "https://gofile.io").
		SetHeader("referer", "https://gofile.io/").
		Get(url_dl)
	if err != nil {
		return "Gofile Changed again", "", err
	}

	defer resp3.Body.Close()

	// very hating moment about why i need to do this
	var response_dl map[string]interface{}
	if err := json.NewDecoder(resp3.Body).Decode(&response_dl); err != nil {
		return "", "", err
	}

	data, ok = response_dl["data"].(map[string]interface{})
	if !ok {
		return "", "", errors.New("data not found in json")
	}

	children, ok := data["children"].(map[string]interface{})
	if !ok {
		return "", "", errors.New("children not found in json")
	}

	var file_name, link_dl string

	for _, child := range children {
		childMap := child.(map[string]interface{})
		file_name = childMap["name"].(string)
		link_dl = childMap["link"].(string)
		break
	}
	return file_name, link_dl, nil

	// body, err := io.ReadAll(resp.Body) // read the response body
	// if err != nil {W
	// 	return "", err // return an empty string and the error
	// }
	// return string(body), err
}

func main() {
	var url string
	fmt.Print("Enter URL: ")
	fmt.Scanln(&url)

	file_name, dlLink, err := gofile(url)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Println("File Name:", file_name)
	fmt.Println("Download Link:", dlLink)
}
