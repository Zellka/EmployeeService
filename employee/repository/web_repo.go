package repository

import (
	"io/ioutil"
	"log"
	"net/http"
)

type WebRepository struct {
	url string
}

func NewWebRepository(url string) *WebRepository {
	return &WebRepository{
		url: url,
	}
}

func (r WebRepository) SetEmployees() []byte {
	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, r.url, nil)
	if err != nil {
		log.Fatal(err)
	}
	res, getErr := client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	return body
}
