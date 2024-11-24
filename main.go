package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pquerna/otp/totp"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cgi"
	"os"
	"path/filepath"
	"strings"
)

func usage() {
	_, _ = fmt.Fprintln(os.Stderr, "TOTP_SECRET=<secret> TOTP_ACCOUNT=<username> mango <port number> <www root directory>")
}

func noopHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	})
}

func logHandler(handler http.Handler) http.Handler {
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

		log.Println(string(d))
	})
}

func main() {
	log.SetOutput(os.Stderr)
	port := "7777"
	root := "."
	AuthenticationHandler := noopHandler

	ServerTOTPsecret, UsingServerTOTP := os.LookupEnv("TOTP_SECRET")
	account, ok := os.LookupEnv("TOTP_ACCOUNT")
	if !ok {
		usage()
		os.Exit(1)
	}

	cgiHandler := func(w http.ResponseWriter, r *http.Request) {
		handler := cgi.Handler{
			Path: filepath.Base(r.URL.Path),
			Dir:  filepath.Join(root, "cgi-bin"),
		}
		handler.ServeHTTP(w, r)
	}
	serverTOTPhandler := func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username, PIN, ok := r.BasicAuth()
			if !ok || strings.TrimSpace(username) != account {
				w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			valid := totp.Validate(PIN, ServerTOTPsecret)
			if !valid {
				http.Error(w, "Unauthorized", 401)
				log.Printf("totp validation error %s\n", PIN)
				return
			}
			handler.ServeHTTP(w, r)
		})
	}

	if UsingServerTOTP {
		AuthenticationHandler = serverTOTPhandler
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

	http.Handle("/", logHandler(AuthenticationHandler(fileServer)))
	http.HandleFunc("/cgi-bin/", logHandler(AuthenticationHandler(http.HandlerFunc(cgiHandler))).ServeHTTP)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
