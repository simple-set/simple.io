package main

import (
	"github.com/simple-set/simple.io/src/example"
	"github.com/sirupsen/logrus"
)

func main() {
	example.NewHttpServerDemo(":8000")

	response, err := example.NewSimpleHttpClient().Get("http://localhost:8000/index?name=x&age=1")
	if err != nil {
		logrus.Errorln(err)
	} else {
		logrus.Infoln(response)
	}
}
