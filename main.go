package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Args[1]
	root := os.Args[2]
	log.Println("Serving from " + root)
	log.Println("Listening on port " + port)
	fs := http.FileServer(http.Dir(root))
	http.Handle("/", fs)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
