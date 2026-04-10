package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

func handleFunc(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: %s \\ %s", r.Method, r.URL.Path)

	if ACLCheck(r.URL.String()) {
		log.Printf("Blocked request to: %s", r.URL.String())
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	log.Printf("\n")

	switch r.Method {
	case http.MethodGet:
		handleGet(w, r)
	case http.MethodConnect:
		handleConnect(w, r)
	default:
		handleAny(w, r)
	}
}

func main() {
	err := loadConfig()
	go loadACL()

	proxy := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
		ReadTimeout:  config.Timeout * time.Second,
		WriteTimeout: config.Timeout * time.Second,
		IdleTimeout:  config.Timeout * time.Second,
		Handler:      securityHeadersMiddleware(http.HandlerFunc(handleFunc)),
	}

	go func() {
		log.Printf("Server started at: %s\n", proxy.Addr)
		err = proxy.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("Initiating graceful shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := proxy.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
