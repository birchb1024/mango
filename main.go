package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cgi"
	"os"
	"path/filepath"
)

func usage() {
	_, _ = fmt.Fprintln(os.Stderr, "mango <port number> <www root directory>")
}

func makeLogHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)

		var buf bytes.Buffer
		tee := io.TeeReader(r.Body, &buf)
		body, _ := ioutil.ReadAll(tee)
		d, _ := json.Marshal(map[string]any{
			"address":        r.RemoteAddr,
			"method":         r.Method,
			"URL":            r.URL,
			"header":         r.Header,
			"content-length": r.ContentLength,
			"body":           string(body),
		})

		log.Println(os.Stderr, string(d))
	})
}

func main() {
	log.SetOutput(os.Stderr)
	port := "7777"
	root := "."

	cgiHandler := func(w http.ResponseWriter, r *http.Request) {
		handler := cgi.Handler{
			Path: filepath.Base(r.URL.Path),
			Dir:  filepath.Join(root, "cgi-bin"),
		}
		handler.ServeHTTP(w, r)
	}

	if len(os.Args) == 3 {
		port = os.Args[1]
		root = os.Args[2]
	} else {
		usage()
		os.Exit(1)
	}
	log.Println("Serving from " + root)
	log.Println("Listening on port " + port)
	fileServer := http.FileServer(http.Dir(root))

	http.Handle("/", makeLogHandler(fileServer))
	http.HandleFunc("/cgi-bin/", makeLogHandler(http.HandlerFunc(cgiHandler)).ServeHTTP)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
