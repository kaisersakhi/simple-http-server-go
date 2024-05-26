package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
)

type transaction func(*Request, *Response)

type route struct {
	method   string
	callback transaction
}

type SimpleServer struct {
	ip           string
	port         string
	routingTable map[string]route
}

func NewServer(ip string, port string) *SimpleServer {
	return &SimpleServer{
		ip:           ip,
		port:         port,
		routingTable: make(map[string]route),
	}
}

func (s *SimpleServer) RegisterRoute(method string, routeRx string, transaction func(request *Request, response *Response)) {
	s.routingTable[routeRx] = route{
		method:   method,
		callback: transaction,
	}
}

func (s *SimpleServer) Register404(transaction func(request *Request, response *Response)) {
	s.routingTable["404"] = route{
		method:   "get",
		callback: transaction,
	}
}

func (s *SimpleServer) Listen() {
	sockAddress := s.ip + ":" + s.port

	listener, err := net.Listen("tcp", sockAddress)

	if err != nil {
		fmt.Printf("Failed on start TCP server: %s", sockAddress)
		os.Exit(1)
	}

	fmt.Printf("Server is listening on: %s", sockAddress)

	defer listener.Close()

	for {
		client, err := listener.Accept()

		if err != nil {
			fmt.Printf("Error accepting connection: %s", err.Error())
			os.Exit(1)
		}

		go handleClient(s, client)
	}
}

func handleClient(server *SimpleServer, client net.Conn) {
	defer client.Close()

	request := PrepareRequest(bufio.NewReader(client))
	response := NewResponse()

	for key, value := range server.routingTable {
		if matchRoute(request.route, key) && strings.EqualFold(value.method, request.method) {
			value.callback(request, response)

			client.Write(response.PrepareRaw())

			return
		}
	}

	// Call 404 callback.

	response.Status(NotFound)

	client.Write(response.PrepareRaw())
}

func matchRoute(route string, rx string) bool {
	regex, err := regexp.Compile(rx)

	if err != nil {
		return false
	}
	return regex.MatchString(route)
}
