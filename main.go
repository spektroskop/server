package main

import (
	"net/http"
	"net/http/httputil"
	"net/http/pprof"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/namsral/flag"
)

var (
	clientDir = flag.String("client", "./client", "")
	listen    = flag.String("listen", ":3000", "")
	proxy     = flag.String("proxy", "", "{prefix},{url};...")
	profile   = flag.String("profile", "", "")
)

func main() {
	flag.Parse()

	router := mux.NewRouter()

	if *proxy != "" {
		proxies := strings.Split(*proxy, ";")

		for _, proxy := range proxies {
			parts := strings.Split(proxy, ",")
			if len(parts) < 2 {
				logrus.Fatalf("Proxy syntax")
			}

			url, err := url.Parse(parts[1])
			if err != nil {
				logrus.Fatalf("Proxy syntax")
			}

			logrus.Infof("Proxy %s -> %s", parts[0], url)

			router.PathPrefix(parts[0]).Handler(
				httputil.NewSingleHostReverseProxy(url),
			)
		}
	}

	if *profile != "" {
		prof := func(name string) string { return filepath.Join("/", *profile, name) }
		router.Handle(prof("profile"), http.HandlerFunc(pprof.Profile))
		router.Handle(prof("heap"), pprof.Handler("heap"))
		router.Handle(prof("trace"), http.HandlerFunc(pprof.Trace))
	}

	router.PathPrefix("/").HandlerFunc(staticHandler(*clientDir)).Methods("HEAD", "GET")

	logrus.Fatal(http.ListenAndServe(*listen, logHandler(router)))
}
