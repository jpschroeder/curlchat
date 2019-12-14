package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	// Accept a command line flag "-httpaddr :8080"
	// This flag tells the server the http address to listen on
	httpaddr := flag.String("httpaddr", "localhost:8080",
		"the address/port to listen on for http \n"+
			"use :<port> to listen on all addresses\n")

	// Accept a command line flag "-baseurl https://mysite.com/"
	baseurl := flag.String("baseurl", "http://localhost:8080/",
		"the base url of the service \n")

	flag.Parse()

	s := Server{&PipeCollection{}, *baseurl}
	http.HandleFunc("/", s.Connect)

	log.Println("Listening on http:", *httpaddr)
	log.Fatal(http.ListenAndServe(*httpaddr, nil))
}
