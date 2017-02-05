package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/trace"

	"github.com/uber-go/zap"
)

func logHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Request",
			zap.String("Method", r.Method),
			zap.String("Path", r.URL.Path),
			zap.String("RemoteAddr", r.RemoteAddr),
		)

		next.ServeHTTP(w, r)
	})
}

func makePath(base, path string) string {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return filepath.Join(base, "index.html")
	} else {
		return path
	}
}

func staticHandler(base string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, makePath(base, filepath.Join(base, r.URL.Path[1:])))
	}
}

func traceHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tr = trace.New(strings.Split(r.URL.Path, "/")[1], r.URL.Path)
		defer tr.Finish()
		tr.LazyPrintf("RemoteAddr: %v", r.RemoteAddr)

		next.ServeHTTP(w, r)
	})
}
