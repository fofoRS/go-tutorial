package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strings"
)

type server struct {
	net.Listener
	handlers map[string]handler
}

func (s *server) RegisterHandle(h handler) error {
	path := h.getPath()
	if path == "" {
		return errors.New("Path is empty, cannot not register the handler")
	}
	s.handlers[path] = h
	return nil
}

func (s *server) start() {
	defer s.Listener.Close()
	fmt.Printf("Starting the server at address %s\n", s.Listener.Addr().String())
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			panic(err)
		}
		firstLine, _ := bufio.NewReader(conn).ReadString('\n')
		path := strings.Fields(firstLine)[1]
		handler, ok := s.handlers[path]
		if ok {
			go handler.handle(conn)
		} else {
			fallbackHandler := s.handlers["/"]
			go fallbackHandler.handle(conn)
		}
	}
}

type handler interface {
	// Mathes the incoming URL path with tha path this handler handles requests.
	// If the path matches, the server will route the request to the handler
	matchPath(string) bool
	getPath() string
	handle(net.Conn)
}

func NewServer(network, address string) server {
	listener, err := net.Listen(network, address)
	if err != nil {
		panic(err)
	}
	return server{Listener: listener, handlers: make(map[string]handler)}
}
