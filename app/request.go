package main

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

type Request struct {
	host          string
	userAgent     string
	contentLength int
	contentType   string
	body          string
	method        string
	httpVersion   string
	route         string
}

func PrepareRequest(reader *bufio.Reader) *Request {
	// Get HTTP method, route, and version.
	parsedheaders := make(map[string]string)

	parseHeaders(parsedheaders, reader)

	contentLength, _ := strconv.Atoi(parsedheaders["Content-Length:"])
	return &Request{
		contentLength: contentLength,
		contentType:   parsedheaders["Content-Type:"],
		body:          parsedheaders["request_body"],
		method:        parsedheaders["method"],
		httpVersion:   parsedheaders["version"],
		route:         parsedheaders["route"],
	}
}

func parseHeaders(parsedHeader map[string]string, reader *bufio.Reader) {
	var bodyEncountered = false

	start_line, _ := reader.ReadString('\n')

	start_line_parts := strings.Split(start_line, " ")

	parsedHeader["method"] = start_line_parts[0]
	parsedHeader["route"] = start_line_parts[1]
	parsedHeader["version"] = start_line_parts[2]

	for {
		if bodyEncountered {
			cLen, isCLenPresent := parsedHeader["Content-Length:"]

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

			parsedHeader["request_body"] = string(buffer)
			fmt.Println("body: ", parsedHeader["request_body"])

			return
		}

		line, err := reader.ReadString('\n')

		if err != nil || line != "\r\n" {
			parts := strings.Split(line, " ")
			parsedHeader[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			fmt.Printf("%v  %v\n", parts[0], parsedHeader[parts[0]])
		} else if line == "\r\n" {
			bodyEncountered = true
		} else {
			break
		}
	}
}
