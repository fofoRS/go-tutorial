package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"runtime/debug"
)

func main() {
	var devMode bool
	flag.BoolVar(&devMode, "dev", false, "tells the app to run in development mode")
	flag.Parse()
	mux := http.NewServeMux()
	mux.HandleFunc("/panic/", panicDemo)
	mux.HandleFunc("/panic-after/", panicAfterDemo)
	mux.HandleFunc("/", hello)
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
					fmt.Fprintf(w, "<h1>%v</h1><pre>%s</pre>", err, stackAsString)
					return
				}
				http.Error(w, "Something went wrong!!", http.StatusInternalServerError)
			}
		}()
		wr := &responseWriterWrapper{ResponseWriter: w}
		mux.ServeHTTP(wr, r)
		wr.flush()
	}
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
	flusher, ok := wrapper.ResponseWriter.(http.Flusher)
	if ok {
		flusher.Flush()
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
		// wrapper.ResponseWriter.WriteHeader(http.StatusOK)
		wrapper.ResponseWriter.Write(b)
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
