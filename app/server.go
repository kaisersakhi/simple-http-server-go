package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
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

	route := headers["route"]

	responseMap := make(map[string]string)
	if matchRoute(route, "/$") {
		fmt.Println("Root route")
		responseMap["status_code"] = "200"

		client.Write(buildResponse(responseMap))
	} else if matchRoute(route, "\\/echo\\/.+") {
		fmt.Println("echo route")
		responseMap["status_code"] = "200"
		responseMap["body"] = strings.SplitAfter(route, "/echo/")[1]

		client.Write(buildResponse(responseMap))
	} else {
		fmt.Println("404 route")
		responseMap["status_code"] = "404"

		client.Write(buildResponse(responseMap))
	}
}

func matchRoute(route string, rx string) bool {
	regex, err := regexp.Compile(rx)

	if err != nil {
		return false
	}
	return regex.MatchString(route)
}

func buildResponse(responseMap map[string]string) []byte{
	var response strings.Builder

	if responseMap["status_code"] == "200" {
		response.WriteString("HTTP/1.1 200 OK\r\n\r\n")
	} else if responseMap["status_code"] == "404" {
		response.WriteString("HTTP/1.1 404 Not Found\r\n\r\n")
	}

	if  responseMap["body"] != "" {
		response.WriteString("Content-Length: " + strconv.Itoa(len([]byte(responseMap["body"]))) + " \r\n\r\n")
		response.WriteString("Content-Type: text/plain\r\n\r\n")
		response.WriteString(responseMap["body"])

	}

//	fmt.Println(string(len(responseMap["body"])))
	fmt.Println(response.String())

	return []byte(response.String())
}
