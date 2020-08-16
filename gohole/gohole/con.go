package gohole

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/textproto"
	"net/url"
	"strings"
)

type holeconn struct {
	rwc    net.Conn
	brc    *bufio.Reader
	server *Server
}

func (hc *holeconn) serve() {
	rawHTTPRequestHeader, remote, credential, isHTTPS, err := hc.getTunnelInfo()
	if err != nil {
		log.Println("Error", err)
		return
	}
	if hc.auth(credential) == false {
		log.Println("Auth fail: " + credential)
		return
	}
	log.Println("connecting to " + remote)
	remoteConn, err := net.Dial("tcp", remote)
	if err != nil {
		log.Println("Error", err)
		return
	}
	if isHTTPS {
		// if https, should sent 200 to client
		_, err = hc.rwc.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
		if err != nil {
			log.Println("Error", err)
			return
		}
	} else {
		// if not https, should sent the request header to remote
		_, err = rawHTTPRequestHeader.WriteTo(remoteConn)
		if err != nil {
			log.Println("Error", err)
			return
		}
	}
	// build bidirectional-streams
	log.Println("begin tunnel", hc.rwc.RemoteAddr(), "<->", remote)
	hc.tunnel(remoteConn)
	log.Println("stop tunnel", hc.rwc.RemoteAddr(), "<->", remote)
}

// getClientInfo parse client request header to get some information:
func (hc *holeconn) getTunnelInfo() (rawReqHeader bytes.Buffer, host, credential string, isHTTPS bool, err error) {
	tp := textproto.NewReader(hc.brc)

	// First line: GET /index.html HTTP/1.0
	var requestLine string
	if requestLine, err = tp.ReadLine(); err != nil {
		return
	}

	method, requestURI, _, ok := parseRequestLine(requestLine)
	if !ok {
		err = &BadRequestError{"malformed HTTP request"}
		return
	}

	// https request
	if method == "CONNECT" {
		isHTTPS = true
		requestURI = "http://" + requestURI
	}

	// get remote host
	uriInfo, err := url.ParseRequestURI(requestURI)
	if err != nil {
		return
	}

	// Subsequent lines: Key: value.
	mimeHeader, err := tp.ReadMIMEHeader()
	if err != nil {
		return
	}

	credential = mimeHeader.Get("Proxy-Authorization")

	if uriInfo.Host == "" {
		host = mimeHeader.Get("Host")
	} else {
		if strings.Index(uriInfo.Host, ":") == -1 {
			host = uriInfo.Host + ":80"
		} else {
			host = uriInfo.Host
		}
	}

	// rebuild http request header
	rawReqHeader.WriteString(requestLine + "\r\n")
	for k, vs := range mimeHeader {
		if k == "Proxy-Authorization" {
			continue
		}
		for _, v := range vs {
			rawReqHeader.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
		}
	}
	rawReqHeader.WriteString("\r\n")
	return
}

// auth provide basic authentication
func (hc *holeconn) auth(credential string) bool {
	if hc.server.isAuth() == false || hc.server.validateCredential(credential) {
		return true
	}
	// 407
	_, err := hc.rwc.Write(
		[]byte("HTTP/1.1 407 Proxy Authentication Required\r\nProxy-Authenticate: Basic realm=\"*\"\r\n\r\n"))
	if err != nil {
		log.Println(err)
	}
	return false
}

// tunnel http message between client and server
func (hc *holeconn) tunnel(remoteConn net.Conn) {
	go func() {
		_, err := hc.brc.WriteTo(remoteConn)
		if err != nil {
			log.Println("Warning", err)
		}
		remoteConn.Close()
	}()
	_, err := io.Copy(hc.rwc, remoteConn)
	if err != nil {
		log.Println("Warning", err)
	}
}
