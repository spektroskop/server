package main

import (
	"encoding/json"
	"net"
	"net/http"
	"net/http/httputil"
	"net/http/pprof"
	"net/url"
	"os"
	"time"

	"golang.org/x/net/trace"

	"github.com/gorilla/mux"
	"github.com/namsral/flag"
	"github.com/uber-go/zap"
)

var (
	clientDir = flag.String("client", "./client", "")
	listen    = flag.String("listen", ":3000", "")
	debug     = flag.Bool("debug", false, "")
	proxies   = flag.String("proxy", "", "")
	logger    = zap.New(
		zap.NewTextEncoder(zap.TextNoTime()),
	)
)

type Proxy map[string][]string

func main() {
	flag.Parse()

	router := mux.NewRouter()

	if *proxies != "" {
		f, err := os.Open(*proxies)
		if err != nil {
			logger.Fatal("Open proxy definition",
				zap.Error(err),
			)
		}

		var proxy Proxy
		if err := json.NewDecoder(f).Decode(&proxy); err != nil {
			logger.Fatal("Decode proxy definition",
				zap.Error(err),
			)
		}

		for endpoint, targets := range proxy {
			for _, target := range targets {
				url, err := url.Parse(endpoint)
				if err != nil {
					logger.Fatal("Endpoint syntax",
						zap.Error(err),
					)
				}

				proxy := httputil.NewSingleHostReverseProxy(url)
				proxy.FlushInterval = 100 * time.Millisecond
				proxy.Transport = &http.Transport{
					Proxy: http.ProxyFromEnvironment,
					Dial: (&net.Dialer{
						Timeout:   24 * time.Hour,
						KeepAlive: 24 * time.Hour,
					}).Dial,
				}

				logger.Info("Define proxy",
					zap.String("Target", target),
					zap.String("Endpoint", url.String()),
				)

				router.PathPrefix(target).Handler(
					proxy,
				)
			}
		}
	}

	var chain http.Handler

	if *debug {
		logger.Info("Debug enabled")

		router.Handle("/debug/profile", http.HandlerFunc(pprof.Profile))
		router.Handle("/debug/trace", http.HandlerFunc(pprof.Trace))
		router.Handle("/debug/heap", pprof.Handler("heap"))
		router.Handle("/debug/requests", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			trace.Render(w, r, false)
		}))

		chain = logHandler(traceHandler(router))
	} else {
		chain = logHandler(router)
	}

	router.PathPrefix("/").HandlerFunc(
		staticHandler(*clientDir),
	).Methods("HEAD", "GET")

	if err := http.ListenAndServe(*listen, chain); err != nil {
		logger.Fatal("Listen", zap.Error(err))
	}
}
