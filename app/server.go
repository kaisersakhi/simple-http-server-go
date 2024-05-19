package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var directory = ""

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	fmt.Printf("Server is listening on port: 4221")

	defer l.Close()

	// Get directory path.
	if len(os.Args) > 2 && os.Args[1] == "--directory" {
		directory = os.Args[2]
	}

	// Wait for clients to join and hand over each client to the goroutine.
	for {
		client, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleClient(client)
	}
}

func handleClient(client net.Conn) {
	defer client.Close()

	// Reads the request data.
	reader := bufio.NewReader(client)

	headers := make(map[string]string)

	// Get the HTTP starting line.
	line, _ := reader.ReadString('\n')

	fmt.Println("Requet headers....")
	fmt.Println(line)

	parts := strings.Split(line, " ")

	// Get individual segments of the first header line.
	headers["action"] = parts[0]  // What is the action? GET, POST, PATCH and etc.
	headers["route"] = parts[1]   // What is the route, ex: /hello-world
	headers["version"] = parts[2] // HTTP version

	parseHeaders(reader, headers)

	route := headers["route"]

	responseMap := make(map[string]string)

	// Decide routes and calculate corresponding response.
	if matchRoute(route, "/$") {
		fmt.Println("Root route")
		responseMap["status_code"] = "200"

		client.Write(buildResponse(responseMap))
	} else if matchRoute(route, "\\/echo\\/.+") {
		fmt.Println("echo route")
		responseMap["status_code"] = "200"
		responseMap["body"] = strings.SplitAfter(route, "/echo/")[1]

		client.Write(buildResponse(responseMap))
	} else if matchRoute(route, "\\/user-agent") {
		fmt.Println("user-agent route")

		responseMap["status_code"] = "200"
		responseMap["body"] = headers["User-Agent:"]

		client.Write(buildResponse(responseMap))
	} else if matchRoute(route, "\\/files\\/.+") {

		if headers["action"] == "POST" {
			fileName := strings.SplitAfter(route, "/files/")[1]

			err := writeContentToFile(directory, fileName, headers["request_body"])

			if err != nil {
				responseMap["status_code"] = "400" // Bad request.
			} else {
				responseMap["status_code"] = "201" // Resource created.
			}

			client.Write(buildResponse(responseMap))
		} else if headers["action"] == "GET" {

			fmt.Println("file path...")

			fileName := strings.SplitAfter(route, "/files/")[1]

			if fileContent, err := fileContentIn(directory, fileName); err != nil {
				responseMap["status_code"] = "404"
			} else {
				responseMap["status_code"] = "200"
				responseMap["body"] = fileContent
				responseMap["content_type"] = "application/octet-stream"
			}
			client.Write(buildResponse(responseMap))
		}
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

func buildResponse(responseMap map[string]string) []byte {
	var response strings.Builder
	var contentType string

	if responseMap["status_code"] == "200" {
		response.WriteString("HTTP/1.1 200 OK\r\n")
	} else if responseMap["status_code"] == "201" {
		response.WriteString("HTTP/1.1 201 OK\r\n")
	} else if responseMap["status_code"] == "404" {
		response.WriteString("HTTP/1.1 404 Not Found\r\n")
	}

	if responseMap["content_type"] != "" {
		contentType = responseMap["content_type"]
	} else {
		contentType = "text/plain"
	}

	if responseMap["body"] != "" {
		response.WriteString("Content-Type: " + contentType + " \r\n")
		response.WriteString("Content-Length: " + strconv.Itoa(len([]byte(responseMap["body"]))) + " \r\n\r\n")
		response.WriteString(responseMap["body"])

	} else {
		response.WriteString("Content-Length: 0 \r\n\r\n")
	}

	fmt.Println(string(len(responseMap["body"])))
	fmt.Println(response.String())

	return []byte(response.String())
}

func parseHeaders(reader *bufio.Reader, headers map[string]string) {
	var bodyEncountered = false

	for {
		if bodyEncountered {
			// reader.ReadString('\n') // Ignore \r\n

			cLen, isCLenPresent := headers["Content-Length:"]

			if !isCLenPresent {
				return
			}

			length, err := strconv.Atoi(cLen)

			if err != nil {
				return
			}

			buffer := make([]byte, length)

			_, err = reader.Read(buffer)

			if err != nil {
				return
			}

			headers["request_body"] = string(buffer)
			fmt.Println("body: ", headers["request_body"])

			return
		}

		line, err := reader.ReadString('\n')

		if err != nil || line != "\r\n" {
			parts := strings.Split(line, " ")
			headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			fmt.Printf("%v  %v\n", parts[0], headers[parts[0]])
		} else if line == "\r\n" {
			bodyEncountered = true
		} else {
			break
		}
	}

}

// Reads a file from the directory specified via --directory flag.
func fileContentIn(directory string, fileName string) (string, error) {
	fileContent, err := os.ReadFile(directory + "/" + fileName)

	if err != nil {
		return "", err
	}

	return string(fileContent), nil
}

// Write content to a file specified.
func writeContentToFile(directory string, fileName string, content string) error {
	absoluteFilePath := filepath.Join(directory, fileName)

	return os.WriteFile(absoluteFilePath, []byte(content), 0644)
}
