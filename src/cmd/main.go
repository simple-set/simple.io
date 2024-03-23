package main

import (
	"github.com/simple-set/simple.io/src/example"
	"time"
)

func main() {
	//server := example.NewRpcServer(":8080")
	//server.Start()

	rpcClient := example.NewRpcClient("localhost:8080")
	rpcClient.Call("hello", []byte("hi"))
	time.Sleep(1000)
}
