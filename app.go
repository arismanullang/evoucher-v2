package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/pkg/profile"
	"github.com/urfave/negroni"

	"github.com/evoucher/voucher/internal/model"
)

var (
	name    = "voucher"
	version = "unversioned"
	token   = name + "/" + version

	fversion = flag.Bool("version", false, "print the version.")
	fconfig  = flag.String("config", "files/etc/voucher/config.yml", "set the config file path.")
	fprofile = flag.String("profile", "", "enable profiler, value either one of [cpu, mem, block].")

	configDir = ""
)

func init() {
	flag.Parse()

	if *fversion {
		printVersion()
	}
}

func main() {
	switch *fprofile {
	case "cpu":
		defer profile.Start(profile.CPUProfile).Stop()
	case "mem":
		defer profile.Start(profile.MemProfile).Stop()
	case "block":
		defer profile.Start(profile.BlockProfile).Stop()
	}

	if err := ReadConfig(*fconfig, &config); err != nil {
		log.Fatal(err)
	}

	if err := model.ConnectDB(config.Database.Endpoint); err != nil {
		log.Fatal(err)
	}

	r := setRoutes()
	m := negroni.New()
	m.Use(negroni.NewRecovery())
	m.UseHandler(r)

	log.Printf("Server is listening on %q\n", config.Server.Host)
	http.ListenAndServe(config.Server.Host, m)
}

func printVersion() {
	fmt.Println(token)
	os.Exit(0)
}
