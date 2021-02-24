package web

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

type HandlerOptionFunc func(h *storyTemplateHTTPHandler)

func withCustomTemple(t *template.Template) HandlerOptionFunc {
	return func(h *storyTemplateHTTPHandler) {
		if t != nil {
			h.Template = t
		}
	}
}

type storyTemplateHTTPHandler struct {
	Story
	*template.Template
}

type Story map[string]Chapter

type Chapter struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []Option `json:"options"`
}

type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

func NewHandler(story Story, handlerOptions ...HandlerOptionFunc) http.Handler {
	defaultTemplate := template.Must(template.New("").Parse("../syos.html"))
	handler := storyTemplateHTTPHandler{story, defaultTemplate}
	for _, option := range handlerOptions {
		option(&handler)
	}
	return handler
}

func (h storyTemplateHTTPHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	path := strings.TrimSpace(request.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}
	key := path[1:]
	chapter, ok := h.Story[key]
	if ok {
		templateExecError := h.Template.Execute(response, chapter)
		if templateExecError != nil {
			log.Fatal(templateExecError)
			http.Error(response, "Error occourred compiling the template.", http.StatusInternalServerError)
		}
		return
	}
	http.Error(response, "Story not found.", http.StatusNotFound)
}

func DecodeJsonFile(r io.Reader) (*Story, error) {
	var story Story
	decoder := json.NewDecoder(r)
	decodeError := decoder.Decode(&story)

	if decodeError != nil {
		return nil, decodeError
	}

	return &story, nil
}
