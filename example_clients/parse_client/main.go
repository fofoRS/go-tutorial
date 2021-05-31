package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fofoRS/go-tutorial/parse"
)

func main() {
	fileNamePointer := flag.String("file", "example.html", "html file used for parse the link found in the file")
	flag.Parse()
	fileNameTokens := strings.Split(*fileNamePointer, ".")
	if len(fileNameTokens) == 1 {
		log.Fatalf("File is not html format %s", *fileNamePointer)
	}
	file, openFileErr := os.Open(*fileNamePointer)
	if openFileErr != nil {
		log.Fatal("Error ocourred opening the file", openFileErr)
	}

	links, err := parse.Parse(file)

	if err != nil {
		log.Fatal("Failed, exiting")
	}

	fmt.Printf("%v\n", links)
}
