package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"runtime"

	"github.com/gilkor/evoucher/lib/x"
	"github.com/gilkor/evoucher/lib/x/jsonerr"

	_ "github.com/joho/godotenv/autoload"
)

func ListenAndServe(h http.Handler) error {
	r := http.NewServeMux()
	r.HandleFunc("/.well-known/acme-challenge/", redirectToSSLManager)
	r.HandleFunc("/debug/ping", ping)
	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)
	r.Handle("/", redirectToHTTPS(h))

	host := os.Getenv("HOST")
	if host == "" {
		host = ":8080"
	}
	log.Printf("Server is listening on port %q\n", host)
	return http.ListenAndServe(host, recoveryHandler(r))
}

var panicHTMLTemplate = template.Must(template.New("_panicHTML").Parse(`
	<!doctype html>
	<html>
	<head>
	<title>Internal Server Error</title>
	</head>

	<body>
		<h1>Internal Server Error</h1>
		<p>An error occurred while processing your request.</p>
	</body>
	</html>
`))

func recoveryHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				stack := make([]byte, 1024*8)
				stack = stack[:runtime.Stack(stack, false)]
				log.Printf("PANIC: %s\n%s", err, stack)

				accept := x.SplitCommaWithTrim(r.Header.Get("Accept"))
				if x.StringInSlice("application/json", accept) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(jsonerr.ErrFatal)
					return
				} else if x.StringInSlice("text/html", accept) {
					w.Header().Set("Content-Type", "text/html")
					w.WriteHeader(http.StatusInternalServerError)
					panicHTMLTemplate.Execute(w, nil)
					return
				} else if x.StringInSlice("text/plain", accept) {
					w.Header().Set("Content-Type", "text/plain")
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, "500 internal server error")
					return
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(jsonerr.ErrFatal)
			}
		}()

		h.ServeHTTP(w, r)
	})
}
