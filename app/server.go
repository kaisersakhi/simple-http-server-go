package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	defer l.Close()

	client, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	handleClient(client)
}

func handleClient(client net.Conn) {
	defer client.Close()

	reader := bufio.NewReader(client)

	headers := make(map[string]string)

	line, _ := reader.ReadString('\n')

	parts := strings.Split(line, " ")

	headers["action"] = parts[0]
	headers["route"] = parts[1]
	headers["version"] = parts[2]

	for {
		line, err := reader.ReadString('\n')

		if err != nil || line != "\r\n" {
			parts = strings.Split(line, " ")
			headers[parts[0]] = parts[1]
			fmt.Printf("%v  %v\n", parts[0], headers[parts[0]])
		} else {
			break
		}
	}

	if headers["route"] == "/" {
		client.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else if headers["route"] == "/echo/abc" {
		client.Write([]byte("HTTP/1.1 200 OK\r\n\r\nabc\r\n"))
	} else {
		client.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
}
