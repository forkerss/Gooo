package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"reflect"
)

var (
	host string
	port int
	dir  string
)

func init() {
	flag.StringVar(&host, "host", "0.0.0.0", "http.server listen host default: \"0.0.0.0\"")
	flag.IntVar(&port, "port", 8080, "http.server listen port default: 8080")
	flag.Parse()
}

func getStatusCode(w http.ResponseWriter) int64 {
	respValue := reflect.ValueOf(w)
	if respValue.Kind() == reflect.Ptr {
		respValue = respValue.Elem()
	}
	status := respValue.FieldByName("status")
	return status.Int()
}

func withLog(h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
		log.Printf("handle %s %d\n", r.URL.Path, getStatusCode(w))
	})
}

func main() {
	addr := fmt.Sprintf("%s:%d", host, port)
	log.Printf("listen on http://%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, withLog(http.FileServer(http.Dir("./")))))
}
