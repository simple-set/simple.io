package main

import (
	"fmt"
	"github.com/simple-set/simple.io/src/event"
	"github.com/simple-set/simple.io/src/handle"
	"github.com/sirupsen/logrus"
)

func server() {
	event.NewBootstrap().
		TcpServer(":8000").
		AddHandler(handle.NewStringDecoder()).
		AddHandler(handle.NewPrintHandler()).
		Bind().Wait()
}

func client() {
	bootstrap := event.NewBootstrap()
	bootstrap.TcpClient("localhost:8000")
	bootstrap.AddHandler(handle.NewStringDecoder())
	bootstrap.AddHandler(handle.NewPrintHandler())
	session := bootstrap.Connect()
	var data string
	for {
		if _, err := fmt.Scanln(&data); err == nil {
			if data == "exit" {
				break
			}
			session.WriteAndFlush(data)
		} else {
			logrus.Errorln(err)
		}
	}
}

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	server()
	//httpServer()
	//client()
}
