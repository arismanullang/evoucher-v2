package main

import (
	"log"
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

	n := negroni.New()
	n.Use(negroni.NewRecovery())
	// n.Use(c.CompanyParamMiddleware())
	n.UseHandler(router)

	// n.UseHandler(testRouter)

	log.Fatal(server.ListenAndServe(n))
}
