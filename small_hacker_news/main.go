package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"
	"sync"
	"time"
)

/*
 1. retrieve all the top stories
 2. interate for each of them
 3. then filter only news storie, checking if it has an url
 4. keep a count how many news has retrieve (up to 30)
 	(should be synchronized) to avoid rece or inconsistencies during the read
 5. keep the original order of each new story
 	( maybe using the index from the slice the ids were stored in)
*/

type NewsCounter struct {
	mu      sync.Mutex
	counter int
}

type WebResponse struct {
	Stories []Story
	Time    time.Duration
}

type StoryItemsIndex struct {
	idx      int
	Story    Story
	ErrrItem error
}

/*
	1. cache lvl1 and cache lv2
	2. cache router
	3. backgound cache manager
*/

func main() {
	var port, numOfStories int
	flag.IntVar(&port, "port", 9092, "Port where the server will listen on.")
	flag.IntVar(&numOfStories, "max_stories", 30, "Max stories to display")
	flag.Parse()

	template := template.Must(template.ParseFiles("index.html"))

	http.HandleFunc("/", dashboardHandler(numOfStories, template))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func dashboardHandler(numOfStories int, tpl *template.Template) func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		startTime := time.Now()
		stories := GetStoryDetails()
		tpl.Execute(response, WebResponse{stories, time.Since(startTime)})
	}
}

func GetStoryDetails() []Story {
	ids, err := GetAllTopStories("https://hacker-news.firebaseio.com")
	if err != nil {
		panic("exiting")
	}
	items := make([]StoryItemsIndex, 0)
	itemChannel := make(chan StoryItemsIndex)
	for i := 0; i < 30; i++ {
		go asyncFetchStory(i, ids[i], itemChannel)
		items = append(items, <-itemChannel)
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].idx < items[j].idx
	})
	stories := make([]Story, 0)
	for _, item := range items {
		stories = append(stories, item.Story)
	}
	return stories
}

func syncFetch(id int) *Story {
	story, err := GetStoryDetail("https://hacker-news.firebaseio.com", id)

	if err != nil {
		panic("exiting")
	}
	return story
}

func asyncFetchStory(idx, id int, storyChan chan StoryItemsIndex) {
	story, err := GetStoryDetail("https://hacker-news.firebaseio.com", id)

	if err != nil {
		storyChan <- StoryItemsIndex{idx, *story, err}
	}
	storyChan <- StoryItemsIndex{idx, *story, err}
}
