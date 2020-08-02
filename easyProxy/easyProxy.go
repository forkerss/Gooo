package main

import (
	"io"
	"log"
	"net"
)

const dstHost = "localhost:5000"

func handle(src net.Conn) {
	// 连接目标 address
	dst, err := net.Dial("tcp", dstHost)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		dst.Close()
		src.Close()
	}()

	// 请求转发
	go func() {
		if _, err := io.Copy(dst, src); err != nil {
			log.Fatalln(err)
		}
	}()
	// 响应转发
	if _, err := io.Copy(src, dst); err != nil {
		log.Fatalln(err)
	}

}
func main() {
	listener, err := net.Listen("tcp", ":20008")
	if err != nil {
		log.Fatalln("Unable to bind to port")
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
			break
		}
		go handle(conn)
	}
}
