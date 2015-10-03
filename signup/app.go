package main

import (
	"flag"
	"fmt"
	"github.com/xtraclabs/signup/handlers"
	"log"
	"net/http"
)

func main() {
	var port = flag.Int("port", -1, "Port to listen on")
	flag.Parse()
	if *port == -1 {
		fmt.Println("Must specify a -port argument")
		return
	}

	log.Println("Listening on port ", *port)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), handlers.Handler())
}
