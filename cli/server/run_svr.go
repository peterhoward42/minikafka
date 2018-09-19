package main

import "fmt"

import "github.com/peterhoward42/toy-kafka/svr"

func main() {
	host := "localhost"
	port := 8086
	svr := svr.NewServer()
	svr.Serve(host, port)
	fmt.Println("Finished")
}
