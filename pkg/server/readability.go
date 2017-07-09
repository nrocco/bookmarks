package server

import (
	"github.com/mauidude/go-readability"
	"io/ioutil"
	"net/http"
)

// FetchContent accepts a valid url as input and returns the sanitised html
func FetchContent(url string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	doc, err := readability.NewDocument(string(body))
	if err != nil {
		return "", err
	}

	return doc.Content(), nil
}
