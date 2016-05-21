package main

import (
	"flag"
	"net/http"
)

func main() {
	var host = flag.String("host", "", "ip address to bind to (the default emtpy means all interfaces)")
	var port = flag.String("port", "32000", "port to bind to")
	flag.Parse()
	router := NewRouter()
	http.ListenAndServe(*host+":"+*port, router)
}
