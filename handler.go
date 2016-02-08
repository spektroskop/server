package main

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/Sirupsen/logrus"
)

func logHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.Infof("%s %s %s", r.RemoteAddr, r.Method, r.RequestURI)
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
