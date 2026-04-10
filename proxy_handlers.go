package main

import (
	"io"
	"net"
	"net/http"
	"sync"
	"time"
)

func handleGet(w http.ResponseWriter, r *http.Request) {
	handleAny(w, r)
}

// NOTE: use curl -x http://127.0.0.1:1234 <url> to test
func handleConnect(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup

	//EXPLINATION: hijacking the connection to handle the CONNECT method
	client_conn, bufrw, err := w.(http.Hijacker).Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	bufrw.Flush()
	defer client_conn.Close()

	// EXPLINATION: establish TCP connection to the target host
	host_conn, err := net.DialTimeout("tcp", r.Host, 3*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer host_conn.Close()

	//EXPLINATION: send 200 Connection Established response to client
	_, err = client_conn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))
	if err != nil {
		return
	}

	// l.Infof("%v -> %v", client_conn.LocalAddr(), host_conn.RemoteAddr())

	//EXPLINATION: start bidirectional piping between client and host
	wg.Go(func() { pipe(client_conn, host_conn) })
	wg.Go(func() { pipe(host_conn, client_conn) })
	wg.Wait() // wait for both goroutines to finish
}

func pipe(src io.Writer, dst io.Reader) {
	_, err := io.Copy(src, dst)
	if err != nil {
		l.Errorf("Error occurred while piping data: %v", err)
	}

	//WARNING: halfclose to prevent race condition
	if t, ok := src.(*net.TCPConn); ok {
		t.CloseWrite()
	}
	if t, ok := dst.(*net.TCPConn); ok {
		t.CloseRead()
	}
}

func handleAny(w http.ResponseWriter, r *http.Request) {

	headers := make(map[string]string)
	for key, values := range r.Header {
		headers[key] = values[0]
	}

	w.Write(MakeRequest(r.URL.String(), r.Method, headers))
}

func MakeRequest(URL string, method string, headers map[string]string) []byte {
	client := &http.Client{}
	req, _ := http.NewRequest(method, URL, nil)

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	removeHopByHopHeaders(req.Header)

	res, err := client.Do(req)
	if err != nil {
		l.Errorf("Err is %v", err)
	}
	defer res.Body.Close()

	resBody, _ := io.ReadAll(res.Body)
	return resBody
}
