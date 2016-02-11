package main

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/namsral/flag"
)

var (
	clientDir = flag.String("client", "./client", "")
	listen    = flag.String("listen", ":3000", "")
)

func main() {
	flag.Parse()

	router := mux.NewRouter()

	router.PathPrefix("/").HandlerFunc(staticHandler(*clientDir)).Methods("HEAD", "GET")

	logrus.Fatal(http.ListenAndServe(*listen, logHandler(router)))
}
