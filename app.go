package main

import (
	"log"
	"os"

	"github.com/gilkor/athena/lib/server"
	"github.com/gilkor/evoucher/internal/model"
	_ "github.com/joho/godotenv/autoload"
	"github.com/urfave/negroni"
)

func main() {
	if err := model.ConnectDB(os.Getenv("DB")); err != nil {
		log.Fatal(err)
	}

	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.UseHandler(router)

	log.Fatal(server.ListenAndServe(n))
}
