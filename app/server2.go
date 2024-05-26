package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var directory string

func main() {
	server := NewServer("0.0.0.0", "4242")
	// Get directory path.
	if len(os.Args) > 2 && os.Args[1] == "--directory" {
		directory = os.Args[2]
	}

	server.RegisterRoute("get", "/$", func(req *Request, res *Response) {
		res.Status(Ok)
		res.ContentType(TextHtml)
		res.WriteBody("<h1 style='color: red;'>Hello world</h1>")
		fmt.Print("Root called..")
	})

	server.RegisterRoute("get", "\\/echo\\/.+", func(req *Request, res *Response) {
		name := strings.SplitAfter(req.route, "/echo/")[1]
		res.WriteBody(name)
	})

	server.RegisterRoute("get", "\\/user-agent", func(req *Request, res *Response) {
		res.WriteBody(req.userAgent)
	})

	server.RegisterRoute("get", "\\/files\\/.+", func(req *Request, res *Response) {
		fileName := strings.SplitAfter(req.route, "/files/")[1]

		fileContent, err := readFileContents(directory, fileName)

		if err != nil {
			res.ResourceNotFound("File requested not found.")
		} else {
			res.WriteBody(fileContent)
		}
	})

	server.RegisterRoute("post", "\\/files\\/.+", func(req *Request, res *Response) {
		fileName := strings.SplitAfter(req.route, "/files/")[1]

		err := writeContentToFile(directory, fileName, req.body)

		if err != nil {
			res.Status(UnprocessableEntity)
		} else {
			res.Status(Created)
		}
	})

	server.Listen()
}

// Reads a file from the directory specified via --directory flag.
func readFileContents(directory string, fileName string) (string, error) {
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
