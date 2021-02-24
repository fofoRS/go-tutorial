package handler

import (
	"fmt"
	"net/http"

	"gopkg.in/yaml.v2"
)

type PathToURL struct {
	path string
	url  string
}

/*
	"MapHandler This methods handles and mappeing the incoming request to be redirect."
*/
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		path := request.URL.Path
		urlRoRedirect, exist := pathsToUrls[path]
		if exist {
			response.Header().Add("Location", urlRoRedirect)
			response.WriteHeader(302)
		} else {
			// pass incoming request to the fallback handler
			fallback.ServeHTTP(response, request)
		}
	}
}

/*
	"YAMLHandler this methods handles and mappeing the incoming request to be redirect."
*/
func YAMLHandler(data []byte, fallback http.Handler) http.HandlerFunc {

	pathToURL, err := parseYaml(data)
	fmt.Printf("%v\n", pathToURL)
	if err != nil {
		fmt.Println(err)
	}
	return MapHandler(pathToURL, fallback)
}

func parseYaml(data []byte) (map[string]string, error) {
	yamlRedirectMapper := make([]map[string]string, 10)
	err := yaml.Unmarshal(data, &yamlRedirectMapper)
	if err != nil {
		return nil, err
	}
	return mapPathToURLAttributes(yamlRedirectMapper), nil

}

func mapPathToURLAttributes(yamlAttributes []map[string]string) map[string]string {
	pathToURLMap := make(map[string]string)
	for _, valueMap := range yamlAttributes {
		tempPathToURL := PathToURL{}
		for key, value := range valueMap {
			if key == "path" {
				tempPathToURL.path = value
			} else {
				tempPathToURL.url = value
			}
		}
		pathToURLMap[tempPathToURL.path] = tempPathToURL.url
	}
	return pathToURLMap
}
