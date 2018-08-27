package main

import "fmt"

import "github.com/peterhoward42/toy-kafka/svr"

func main() {
    svr := svr.NewServer(42)
    svr.Run("fibble")
	fmt.Println("Finished")
}
