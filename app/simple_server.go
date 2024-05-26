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
	routeRx string
	method  string
}

type SimpleServer struct {
	ip           string
	port         string
	routingTable map[route]transaction
}

func NewServer(ip string, port string) *SimpleServer {
	return &SimpleServer{
		ip:           ip,
		port:         port,
		routingTable: make(map[route]transaction),
	}
}

func (s *SimpleServer) RegisterRoute(method string, routeRx string, transaction func(request *Request, response *Response)) {
	route := route{
		routeRx: routeRx,
		method:  method,
	}
	s.routingTable[route] = transaction
}

func (s *SimpleServer) Register404(transaction func(request *Request, response *Response)) {
	route := route{
		routeRx: "\\/404",
		method:  "get",
	}
	s.routingTable[route] = transaction
}

func (s *SimpleServer) Listen() {
	sockAddress := s.ip + ":" + s.port

	listener, err := net.Listen("tcp", sockAddress)

	if err != nil {
		fmt.Printf("Failed on start TCP server: %s %s", sockAddress, err)
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

	for kRoute, vTransaction := range server.routingTable {
		if matchRoute(request.route, kRoute.routeRx) && strings.EqualFold(kRoute.method, request.method) {
			vTransaction(request, response)

			client.Write(response.PrepareRaw())

			return
		}
	}

	// Call 404 callback.

	response.ResourceNotFound("The resource you're looking for is not found.", "Requested resource : "+request.route)

	client.Write(response.PrepareRaw())
}

func matchRoute(route string, rx string) bool {
	regex, err := regexp.Compile(rx)

	if err != nil {
		return false
	}
	return regex.MatchString(route)
}
