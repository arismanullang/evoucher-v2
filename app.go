package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gilkor/athena/lib/server"

	// c "github.com/gilkor/evoucher-v2/controller"
	"github.com/gilkor/evoucher-v2/model"
	_ "github.com/joho/godotenv/autoload"
	"github.com/urfave/negroni"
)

func main() {
	if err := model.ConnectDB(os.Getenv("DB")); err != nil {
		log.Fatal(err)
	}

	model.RegisterValidator()

	n := negroni.New()
	n.UseFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		if r.Method == http.MethodOptions {
			w.Header().Add("Access-Control-Allow-Headers", "Authorization, Accept, Content-Type")
			w.Header().Add("Access-Control-Allow-Methods", "GET, HEAD, POST, PUT, PATCH, DELETE")
			return
		}
		next(w, r)
	})
	n.Use(negroni.NewRecovery())
	// n.Use(c.CompanyParamMiddleware())
	n.UseHandler(router)

	// n.UseHandler(testRouter)

	log.Fatal(server.ListenAndServe(n))
}
