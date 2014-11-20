package main

import (
	"flag"
	"github.com/elazarl/goproxy"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

func GetIP(r *http.Request) string {
	if ipProxy := r.Header.Get("X-FORWARDED-FOR"); len(ipProxy) > 0 {
		return ipProxy
	}

	return r.RemoteAddr
}

func main() {
	verbose := flag.Bool("v", false, "print verbose output")
	addr := flag.String("a", ":8080", "proxy listen address")
	pidfile := flag.String("p", "", "path to pidfile")
	logfile := flag.String("l", "", "path to writable logfile")

	flag.Parse()

	// Write log file.
	var logger *log.Logger
	if *logfile != "" {
		log_f, err := os.OpenFile(*logfile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
		if err == nil {
			defer log_f.Close()
			logger = log.New(log_f, "proxyd: ", log.LstdFlags)
		} else {
			log.Printf("Failed to open log file %s : %v", log_f, err)
			os.Exit(1)
		}
	} else {
		logger = log.New(os.Stdout, "proxyd: ", log.LstdFlags)
	}

	// Write pid file.
	if *pidfile != "" {
		pid := strconv.Itoa(os.Getpid())
		if err := ioutil.WriteFile(*pidfile, []byte(pid), 0644); err != nil {
			log.Printf("Failed to write to pid file : %v", err)
			os.Exit(1)
		}
	}

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = *verbose

	// Log requests
	proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		logger.Printf("Request: %s | IP: %s", req.URL, GetIP(req))
		// If the HTTP response is not nil, the proxy will never send the request to the remote client
		return req, nil
	})

	log.Fatal(http.ListenAndServe(*addr, proxy))
}
