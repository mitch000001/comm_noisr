// +build !appengine
package main

import (
	"log"
	"net/http"
	"os"
)

const DefaultPort string = "5000"

var httpAddress string

func init() {
	port := os.Getenv("PORT")
	if port == "" {
		port = DefaultPort
	}
	httpAddress = ":" + port
}

func main() {
	log.Printf("listening on %s", httpAddress)
	log.Fatal(http.ListenAndServe(httpAddress, http.Handler(router)))
}
