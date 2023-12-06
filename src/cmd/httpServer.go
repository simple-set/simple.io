package main

import (
	"fmt"
	"log"
	"net/http"
)

func indexFunc(response http.ResponseWriter, request *http.Request) {
	body := make([]byte, request.ContentLength)
	request.Body.Read(body)
	fmt.Println(string(body))
	_, err := response.Write([]byte("hello"))
	if err != nil {
		panic(err)
	}
	log.Printf("host: %s, method: %s, path: %s, code: %d", request.RemoteAddr, request.Method, request.RequestURI, 200)
}

func httpServer() {
	http.HandleFunc("/index", indexFunc)
	http.ListenAndServe("localhost:8080", nil)
}
