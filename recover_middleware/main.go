package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

func main() {
	var devMode bool
	flag.BoolVar(&devMode, "dev", true, "tells the app to run in development mode")
	flag.Parse()
	mux := http.NewServeMux()
	mux.HandleFunc("/panic/", panicDemo)
	mux.HandleFunc("/panic-after/", panicAfterDemo)
	mux.HandleFunc(`/debug/`, renderSourceFile)
	mux.HandleFunc("/hello", hello)
	log.Fatal(http.ListenAndServe(":8080", middleWareRecover(mux, devMode)))
}

func middleWareRecover(mux http.Handler, dev bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				stack := debug.Stack()
				stackAsString := string(stack)
				log.Printf(stackAsString)
				if dev {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprint(w, "<html><body>")
					for _, path := range parseStackTrace(stackAsString) {
						fmt.Fprintf(w, path)
					}
					fmt.Fprint(w, "</body></html>")
					return
				}
				http.Error(w, "Something went wrong!!", http.StatusInternalServerError)
			}
		}()
		mux.ServeHTTP(w, r)
	}
}

func parseStackTrace(stack string) []string {
	lines := strings.Split(stack, "\n")
	parsed := make([]string, 0)
	for _, line := range lines {
		if strings.HasPrefix(line, "\t") && strings.Contains(line, ":") {
			// parsed = append(parsed, strings.TrimSpace(line[:strings.Index(line, ":")]))
			path := line[:strings.Index(line, ":")]
			fields := strings.Split(line[strings.Index(line, ":")+1:], " ")
			var lineNumber int
			if len(fields) >= 1 {
				lineNumber, _ = strconv.Atoi(fields[0])
			}
			parsed = append(parsed, fmt.Sprintf(`<a href="/debug/?line=%d&path=%s"><pre>%s:%d</pre></a>`, lineNumber, path, path, lineNumber))
		}
		parsed = append(parsed, line)
	}
	return parsed
}

type responseWriterWrapper struct {
	http.ResponseWriter
	status int
	data   [][]byte
}

func (wrapper *responseWriterWrapper) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj, ok := wrapper.ResponseWriter.(http.Hijacker) // type assertion
	if !ok {
		panic("Response writer doesn't implment Hijacker interface")
	}
	return hj.Hijack()
}

func (wrapper *responseWriterWrapper) Flush() {
	_, ok := wrapper.ResponseWriter.(http.Flusher)
	if ok {
		if wrapper.status != 0 {
			wrapper.ResponseWriter.WriteHeader(wrapper.status)
		}
		for _, b := range wrapper.data {
			wrapper.ResponseWriter.Write(b)
		}
	}
}

func (wrapper *responseWriterWrapper) Write(b []byte) (int, error) {
	wrapper.data = append(wrapper.data, b)
	return len(b), nil
}

func (wrapper *responseWriterWrapper) WriteHeader(statusCode int) {
	wrapper.status = statusCode
}

func (wrapper *responseWriterWrapper) flush() {
	if wrapper.status != 0 {
		wrapper.ResponseWriter.WriteHeader(wrapper.status)
	}
	for _, b := range wrapper.data {
		wrapper.ResponseWriter.Write(b)
	}
}

func renderSourceFile(w http.ResponseWriter, r *http.Request) {
	path := r.FormValue("path")
	line, _ := strconv.Atoi(r.FormValue("line"))
	file, err := os.Open(path)
	if err != nil {
		http.Error(w, "Error reading source file", http.StatusInternalServerError)
	}
	b := bytes.NewBuffer(nil)
	_, err = io.Copy(b, file)
	if err != nil {
		http.Error(w, "Error copying source file", http.StatusInternalServerError)
	}

	lexer := lexers.Get("go")
	style := styles.Get("github")
	if style == nil {
		style = styles.Fallback
	}
	iterator, err := lexer.Tokenise(nil, b.String())
	var lines [][2]int
	if line > 0 {
		lines = append(lines, [2]int{line, line})
	}
	w.Header().Set("Content-Type", "text/html")
	htmlFormatter := html.New(html.TabWidth(2), html.WithLineNumbers(true), html.HighlightLines(lines))
	err = htmlFormatter.Format(w, style, iterator)
	if err != nil {
		panic(err)
	}
}

func panicDemo(w http.ResponseWriter, r *http.Request) {
	funcThatPanics()
}

func panicAfterDemo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<pre>Hello!</pre>")
	funcThatPanics()
}

func funcThatPanics() {
	panic("Oh no!")
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>Hello!</h1>")
}
