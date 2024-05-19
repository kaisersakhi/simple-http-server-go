package main

import "fmt"

func main() {
	server := NewServer("0.0.0.0", "4242")

	server.RegisterRoute("get", "/$", func(req *Request, res *Response) {
		res.Status(Ok)
		res.ContentType(TextHtml)
		res.WriteBody("<h>Hello world</h1>")
		fmt.Print("Root called..")
	})

	server.Listen()
}
