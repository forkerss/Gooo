package main

import (
	"flag"
	"gohole/gohole"
	"log"
)

func init() {
	log.SetPrefix("GoHole: ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
}
func main() {
	http := flag.String("http", ":1081", "proxy listen addr")
	auth := flag.String("auth", "", "basic credentials(username:password)")
	gen := flag.Bool("gen", false, "has gen basic credentials(username:password)")
	flag.Parse()
	server := gohole.NewServer(*http, *auth, *gen)
	server.Start()
}
