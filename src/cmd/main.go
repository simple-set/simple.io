package main

import "github.com/simple-set/simple.io/src/example"

func main() {
	new(example.SimpleHttpServer).Start()
	//new(example.SimpleHttpClient).Connect()
}
