package server

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	. "github.com/leon332157/replish/common"
	log "github.com/sirupsen/logrus"
)

// TODO: Handler for __dav, *.git, __ws, __ssh and wildcard (reverse proxy)
func StartMain(config *ReplishConfig) {

	if config.LocalAppPort == 0 {
		log.Fatal("app port is 0")
	}
	/*http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", r.URL.Path)
	})

	http.HandleFunc("/__dav", handlerDav)
	*/
	//http.FileServer(http.Dir("/home/runner/replish"))
	listener, err := net.Listen("tcp4", fmt.Sprintf(":%v", config.ListenPort))
	if err != nil {
		log.Panicf("[Server Main] %s\n", err)
	}
	log.Infof("[Server Main] Listening on %v", config.ListenPort)
	// p := &ReverseProxy{port: port}
	http.Serve(listener, &ReplishRouter{port: config.LocalAppPort})
	/*go func() {=
		err := http.Serve(listener, p)

		if err != nil {
			log.Panicf("[Server Main] %s\n", err)
		}
	}()

	if err != nil {
		log.Panicf("[Server Main] %s\n", err)
	}
	log.Debug("[Server Main] reverse proxy started")*/
}

type ReplishRouter struct {
	port uint16
}

func (s *ReplishRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if strings.HasPrefix(path, "/__dav") {
		log.Debug("[Server Router] Match /__dav, passing to webdav")
		handlerDav(w, r)
	} else if strings.HasPrefix(path, "/__ws") {
		log.Debug("[Server Router] Matching /__ws, passing to websocket")
		handleWS(w, r)
	} else {
		localUrl, err := url.Parse(fmt.Sprintf("http://127.0.0.1:%v", s.port))
		if err != nil {
			log.Fatalf("[Server Router] Formatting url failed!")
		}
		proxy := httputil.NewSingleHostReverseProxy(localUrl)
		proxy.ServeHTTP(w, r)
	}
}
