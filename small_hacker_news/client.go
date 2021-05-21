package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Story struct {
	Title string `json:"title"`
	Url   string `json:"url"`
}

func GetAllTopStories(baseURL string) ([]int, error) {
	fullURL := fmt.Sprintf("%s/v0/topstories.json", baseURL)
	response, getError := http.Get(fullURL)
	if getError != nil {
		return nil, getError
	}
	defer response.Body.Close()
	decoder := json.NewDecoder(response.Body)
	ids := make([]int, 0)
	decodeErr := decoder.Decode(&ids)
	if decodeErr != nil {
		return nil, decodeErr
	}
	return ids, nil
}

func GetStoryDetail(baseURL string, id int) (*Story, error) {
	fullURL := fmt.Sprintf("%s/v0/item/%d.json", baseURL, id)
	response, getError := http.Get(fullURL)
	if getError != nil {
		return nil, getError
	}
	defer response.Body.Close()
	decoder := json.NewDecoder(response.Body)
	story := &Story{}
	decodeErr := decoder.Decode(story)
	if decodeErr != nil {
		return nil, decodeErr
	}
	return story, nil
}
