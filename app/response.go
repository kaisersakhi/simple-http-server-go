package main

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	Continue            = 100
	Ok                  = 200
	Created             = 201
	BadRequest          = 400
	Unauthorized        = 401
	NotFound            = 404
	UnprocessableEntity = 422
)

const (
	TextPlain              = "text/plain"
	TextHtml               = "text/html"
	ApplicationOctetStream = "application/octet-stream"
)

type Response struct {
	httpVersion   string
	date          string
	contentType   string
	contentLength int
	body          string
	status        int
}

func NewResponse() *Response {
	// Return response object with default values.
	return &Response{
		httpVersion:   "1.1",
		status:        Ok,
		contentType:   TextPlain,
		date:          "",
		contentLength: 0,
		body:          "",
	}
}

func (r *Response) WriteBody(content string) {
	r.body = content
	r.contentLength = len(content)
}

func (r *Response) Status(status int) {
	r.status = status
}

func (r *Response) ContentType(contentType string) {
	r.contentType = contentType
}

func (r *Response) ResourceNotFound(reasons ...string) {
	if len(reasons) == 0 {
		r.WriteBody("Resource not found.")
	} else {
		r.WriteBody(strings.Join(reasons, "\n"))
	}
	r.status = NotFound
}

func (r *Response) PrepareRaw() []byte {
	var rawResponse strings.Builder
	responseMap := make(map[string]string)

	statusLine := "HTTP/" + r.httpVersion + " " + strconv.Itoa(r.status) + " " + reasonPhrase(r.status)

	responseMap["Date:"] = r.date
	responseMap["Content-Type:"] = r.contentType
	responseMap["Content-Length:"] = strconv.Itoa(r.contentLength)

	rawResponse.WriteString(statusLine + "\r\n")

	for key, value := range responseMap {
		if value != "" {
			rawResponse.WriteString(key + " " + value + "\r\n")
		}
	}

	rawResponse.WriteString("\r\n")

	if r.body != "" {
		rawResponse.WriteString(r.body)
	}

	str := rawResponse.String()

	fmt.Println(str)

	return []byte(str)
}

func reasonPhrase(statusCode int) string {
	switch statusCode {
	case 100:
		return "Continue"
	case 200:
		return "OK"
	case 201:
		return "Created"
	case 400:
		return "Bad Request"
	case 401:
		return "Unauthorized"
	case 404:
		return "Not Found"
	case 422:
		return "Unprocessable Entity"
	default:
		return ""
	}
}
