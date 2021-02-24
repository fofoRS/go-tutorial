package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"example.com/story_web/web"
)

func main() {
	file, err := os.Open("story.json")
	if err != nil {
		panic(err)
	}
	story, decodeError := web.DecodeJsonFile(file)
	if decodeError != nil {
		panic(decodeError)
	}
	storyHTMLTemplateBytes, fileReadError := ioutil.ReadFile("syos.html")
	if fileReadError != nil {
		panic(fileReadError)
	}
	storyTemplate := string(storyHTMLTemplateBytes)
	handler := web.NewHandler(story, storyTemplate)
	fmt.Println("starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
