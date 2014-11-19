package main

import (
	"github.com/elazarl/goproxy"
	"log"
	"flag"
	"net/http"
	"os"
	"io/ioutil"
	"strconv"
)

func main() {
	verbose := flag.Bool("v", false, "should every proxy request be logged to stdout")
	addr := flag.String("addr", ":8080", "proxy listen address")
	pidFile := flag.String("pid", "", "the pid file")
        flag.Parse()
        
	// Write pid file.
	if *pidFile != "" {
		pid := strconv.Itoa(os.Getpid())
		if err := ioutil.WriteFile(*pidFile, []byte(pid), 0644); err != nil {
			panic(err)
		}
	}
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = *verbose
	log.Fatal(http.ListenAndServe(*addr, proxy))
}
