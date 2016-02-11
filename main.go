package main

import (
	"flag"
	"net/http"
	"net/http/pprof"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

var (
	clientDir = flag.String("client", "./client", "")
	listen    = flag.String("listen", ":3000", "")
	profile   = flag.Bool("profile", false, "")
)

func main() {
	flag.Parse()

	router := mux.NewRouter()

	if *profile {
		router.Handle("/debug/profile", http.HandlerFunc(pprof.Profile))
		router.Handle("/debug/heap", pprof.Handler("heap"))
		router.Handle("/debug/trace", http.HandlerFunc(pprof.Trace))
	}

	router.PathPrefix("/").HandlerFunc(staticHandler(*clientDir)).Methods("HEAD", "GET")

	logrus.Fatal(http.ListenAndServe(*listen, logHandler(router)))
}
