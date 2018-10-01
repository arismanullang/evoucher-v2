package server

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
)

func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

func redirectToSSLManager(w http.ResponseWriter, r *http.Request) {
	u := &url.URL{}
	*u = *r.URL
	u.Scheme = "http"
	u.Host = "ssl-manager.apps.id"
	http.Redirect(w, r, u.String(), http.StatusFound)
}

func redirectToHTTPS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if os.Getenv("FORCE_HTTPS") == "true" && r.Header.Get("X-Forwarded-Proto") == "http" {
			u := &url.URL{}
			*u = *r.URL
			u.Scheme = "https"
			u.Host = r.Header.Get("Host")
			http.Redirect(w, r, u.String(), http.StatusMovedPermanently)
			return
		}
		h.ServeHTTP(w, r)
	})
}
