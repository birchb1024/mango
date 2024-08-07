package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func usage() {
	_, _ = fmt.Fprintln(os.Stderr, "mango <port number> <www root directory>")
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := new(strings.Builder)
		_, _ = io.Copy(buf, r.Body)
		d, _ := json.Marshal(map[string]any{
			"address":        r.RemoteAddr,
			"method":         r.Method,
			"URL":            r.URL,
			"header":         r.Header,
			"content-length": r.ContentLength,
			"body":           buf.String(),
		})
		log.Println(string(d))
		handler.ServeHTTP(w, r)
	})
}

func main() {
	port := "7777"
	root := "."
	if len(os.Args) == 3 {
		port = os.Args[1]
		root = os.Args[2]
	} else {
		usage()
		os.Exit(1)
	}
	log.Println("Serving from " + root)
	log.Println("Listening on port " + port)
	fs := http.FileServer(http.Dir(root))
	http.Handle("/", fs)

	err := http.ListenAndServe(":"+port, logRequest(http.DefaultServeMux))
	if err != nil {
		log.Fatal(err)
	}
}
