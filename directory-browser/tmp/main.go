package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	args := os.Args[1:]
	source := ""
	cert := ""
	key := ""
	if len(args) > 1 {
		source = args[0]
		cert = args[1]
		key = args[2]
	} else {
		source = "tmp"
		cert = "blog.d0zingcat.xyz.pem"
		key = "blog.d0zingcat.xyz.key"
	}
	http.Handle("/fileserver/", http.StripPrefix("/fileserver/", http.FileServer(http.Dir(source))))
	http.ListenAndServe(":9000", nil)
	http.ListenAndServeTLS(":9000", cert, key, nil)
}
