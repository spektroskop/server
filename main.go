package main

import (
	"flag"
	"net/http"
	"net/http/pprof"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

var (
	clientDir = flag.String("client", "./client", "")
	listen    = flag.String("listen", ":3000", "")
	profile   = flag.String("profile", "", "")
)

func main() {
	flag.Parse()

	router := mux.NewRouter()

	if *profile != "" {
		prof := func(name string) string { return filepath.Join("/", *profile, name) }
		router.Handle(prof("profile"), http.HandlerFunc(pprof.Profile))
		router.Handle(prof("heap"), pprof.Handler("heap"))
		router.Handle(prof("trace"), http.HandlerFunc(pprof.Trace))
	}

	router.PathPrefix("/").HandlerFunc(staticHandler(*clientDir)).Methods("HEAD", "GET")

	logrus.Fatal(http.ListenAndServe(*listen, logHandler(router)))
}
