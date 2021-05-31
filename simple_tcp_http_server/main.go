package main

import (
	"fmt"
	"net"
)

type echoHandler struct {
	path string
}

func (h *echoHandler) matchPath(p string) bool {
	if h.path == p {
		return true
	}
	return false
}

func (h *echoHandler) getPath() string {
	return h.path
}

func (h *echoHandler) handle(conn net.Conn) {
	defer conn.Close()
	fmt.Fprint(conn, "Fallback handler")
}

type htmlHandler struct {
	path string
}

func (h *htmlHandler) matchPath(p string) bool {
	if h.path == p {
		return true
	}
	return false
}

func (h *htmlHandler) getPath() string {
	return h.path
}

func (h *htmlHandler) handle(conn net.Conn) {
	defer conn.Close()
	html := `
		<html>
			<head>
  				<meta charset="utf-8">
			</head>
			<body>
				<h1>Welcome!</h1>
			</body>
		</html>
	`
	fmt.Fprintln(conn, "HTTP/1.1 OK 200")
	fmt.Fprintln(conn, "Connection: close")
	fmt.Fprintln(conn, "Content-Type: text/html")
	fmt.Fprintln(conn, fmt.Sprintf("Content-Length: %d", len(html)))
	fmt.Fprintf(conn, "\n")
	fmt.Fprintf(conn, html)
}

func main() {
	server := NewServer("tcp", "127.0.0.1:8080")
	_ = server.RegisterHandle(&echoHandler{path: "/"})
	_ = server.RegisterHandle(&htmlHandler{path: "/html"})
	server.start()
}
