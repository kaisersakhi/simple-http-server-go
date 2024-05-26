package main

import (
	"fmt"
	"strings"
)

func main() {
	server := NewServer("0.0.0.0", "4242")

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
		res.ResourceNotFound()
	})

	server.RegisterRoute("post", "\\/files\\/.+", func(req *Request, res *Response) {

	})

	server.Listen()
}
