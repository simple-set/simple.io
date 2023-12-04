package main

import (
	"bufio"
	"github.com/simple-set/simple.io/src/event"
	"github.com/simple-set/simple.io/src/handle"
	"os"
)

func httpClient() {
	//var request string = "POST /index HTTP/1.1\nHost: localhost:8080\nConnection: keep-alive\nContent-Length: 19\nSec-Ch-Ua: \"Google Chrome\";v=\"119\", \"Chromium\";v=\"119\", \"Not?A_Brand\";v=\"24\"\nContent-Type: text/plain;charset=UTF-8\nCache-Control: no-cache\nSec-Ch-Ua-Mobile: ?0\nUser-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\nSec-Ch-Ua-Platform: \"Windows\"\nAccept: */*\nOrigin: chrome-extension://coohjcphdfgbiolnekdpbcijmhambjff\nSec-Fetch-Site: none\nSec-Fetch-Mode: cors\nSec-Fetch-Dest: empty\nAccept-Encoding: gzip, deflate, br\nAccept-Language: zh-CN,zh;q=0.9,en;q=0.8\n\n{\n\n  \"name\": \"xk\"\n}"
	var request1 string = "POST /index HTTP/1.1\nHost: localhost:8080\nConnection: keep-alive\nContent-Length: 19\nSec-Ch-Ua: \"Google Chrome\";v=\"119\", \"Chromium\";v=\"119\", \"Not?A_Brand\";v=\"24\"\nContent-Type: text/plain;charset=UTF-8\nCache-Control: no-cache\nSec-Ch-Ua-Mobile: ?0\nUser-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36\nSec-Ch-Ua-Platform: \"Windows\"\nAccept: */*\nOrigin: chrome-extension://coohjcphdfgbiolnekdpbcijmhambjff"
	var request2 string = "Sec-Fetch-Site: none\nSec-Fetch-Mode: cors\nSec-Fetch-Dest: empty\nAccept-Encoding: gzip, deflate, br\nAccept-Language: zh-CN,zh;q=0.9,en;q=0.8\n\n{\n\n  \"name\": \"xk\"\n}"

	bootstrap := event.NewBootstrap()
	bootstrap.TcpClient("localhost:8080")
	bootstrap.AddHandler(handle.NewStringDecoder())
	bootstrap.AddHandler(handle.NewPrintHandler())
	session := bootstrap.Connect()

	for {
		scanner := bufio.NewScanner(os.Stdin)
		for {
			// 一个请求分为两次发送, 模拟半包效果
			scanner.Scan()
			session.WriteAndFlush(request1)
			scanner.Scan()
			session.WriteAndFlush(request2)
		}
	}
}
