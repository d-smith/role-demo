package main

import (
	"flag"
	"fmt"
	"github.com/xtraclabs/signup/handlers"
	"log"
	"net/http"
	"os"
)

func validateEnvironment() {
	if os.Getenv("ROLL_ENDPOINT") == "" {
		println("ROLL_ENDPOINT environment variable not set.")
		os.Exit(1)
	}
}

func main() {

	validateEnvironment()

	var port = flag.Int("port", -1, "Port to listen on")
	flag.Parse()
	if *port == -1 {
		fmt.Println("Must specify a -port argument")
		return
	}

	log.Println("Listening on port ", *port)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), handlers.Handler())
}
