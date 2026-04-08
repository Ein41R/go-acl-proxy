package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// WARNING: incomplete list of hop-by-hop headers
// which will be stripped
var perHopHeaders = []string{
	"Proxy-Connection",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Connection",
	"Keep-Alive",
	"TE",
	"Trailer",
	"Transfer-Encoding",
	"Upgrade",
}

func main() {
	err := loadConfig()
	go loadACL()

	proxy := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
		ReadTimeout:  config.Timeout * time.Second,
		WriteTimeout: config.Timeout * time.Second,
		IdleTimeout:  config.Timeout * time.Second,
		Handler:      http.HandlerFunc(handleFunc),
	}

	host := config.Host
	port := config.Port

	log.Printf("Server started at %s:%d\n", host, port)
	err = proxy.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func handleFunc(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: %s %s\n", r.Method, r.URL.Path)

	if ACLCheck(r.URL.String()) {
		log.Printf("Blocked request to: %s\n", r.URL.String())
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	switch r.Method {
	case http.MethodGet:
		handleGet(w, r)
	case http.MethodConnect:
		handleConnect(w, r)
	default:
		handleAny(w, r)
	}
}
