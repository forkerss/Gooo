package gohole

import (
	"bufio"
	"encoding/base64"
	"log"
	"net"
	"strings"
)

// BadRequestError proxy handle bad request error
type BadRequestError struct {
	what string
}

func (b *BadRequestError) Error() string {
	return b.what
}

// Server proxy server
type Server struct {
	addr       string
	credential string
	listener   net.Listener
}

// NewServer create a proxy server
func NewServer(Addr, credential string, gen bool) *Server {
	if gen {
		credential = RandStringBytesMaskImprSrc(16) + ":" +
			RandStringBytesMaskImprSrc(16)
	}
	return &Server{addr: Addr, credential: credential}
}

// Start proxy server to hanle conn
func (s *Server) Start() {
	var err error
	s.listener, err = net.Listen("tcp", s.addr)
	if err != nil {
		log.Fatalln(err)
	}
	if s.credential != "" {
		log.Printf("use %s for auth\n", s.credential)
	}
	log.Printf("proxy listen in %s, waiting for connection...\n", s.addr)
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Println("ERROR", err)
			continue
		}
		go s.newConn(conn).serve()
	}
}

// newConn create a conn to serve client request
func (s *Server) newConn(rwc net.Conn) *holeconn {
	return &holeconn{
		server: s,
		rwc:    rwc,
		brc:    bufio.NewReader(rwc),
	}
}

func (s *Server) isAuth() bool {
	return s.credential != ""
}

func (s *Server) validateCredential(basicCredential string) bool {
	c := strings.Split(basicCredential, " ")
	if len(c) == 2 && strings.EqualFold(c[0], "Basic") {
		decoded, err := base64.StdEncoding.DecodeString(c[1])
		if err != nil {
			log.Println("Error", err)
			return false
		}
		if s.credential == string(decoded) {
			return true
		}
	}
	return false
}
