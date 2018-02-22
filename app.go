package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"google.golang.org/appengine"

	"github.com/pkg/profile"
	"github.com/ruizu/render"
	"github.com/urfave/negroni"

	"github.com/gilkor/evoucher/internal/controller"
	"github.com/gilkor/evoucher/internal/model"
)

//var Session *redis.Client

var (
	name        = "voucher"
	version     = "unversioned"
	token       = name + "/" + version
	path_config = os.Getenv("EVOUCHER_CONFIG")

	fversion = flag.Bool("version", false, "print the version.")
	// fconfig  = flag.String("config", path_config, "set the config file path.")
	fconfig  = flag.String("config", "files/etc/voucher/config.yml", "set the config file path.")
	fprofile = flag.String("profile", "", "enable profiler, value either one of [cpu, mem, block].")

	configDir = ""
)

func init() {

	flag.Parse()

	if *fversion {
		printVersion()
	}

	/*
		// init redist
		Session = redis.NewClient(&redis.Options{
			Addr:         ":8889",
			DialTimeout:  10 * time.Second,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			PoolSize:     10,
			PoolTimeout:  30 * time.Second,
		})
		Session.FlushDb()
	*/
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
	if err := model.OpenRedisPool(config.Database.Redis.Endpoint); err != nil {
		log.Fatal(err)
	}

	render.SetPath(config.Server.TemplateDirectory)

	model.UiFeatures = getUiRole()
	model.ApiFeatures = getApiRole()
	model.Config = getConfig()
	model.Domain = config.Mailgun.Domain
	model.ApiKey = config.Mailgun.MailgunKey
	model.PublicApiKey = config.Mailgun.MailgunPublicKey
	model.RootTemplate = config.Mailgun.RootTemplate
	model.Email = config.Mailgun.Email
	model.RootUrl = config.Mailgun.RootUrl

	model.GetProgramTypes()

	//logger config
	model.Path = config.Logger.Path
	model.FileName = config.Logger.FileName
	//voucher config
	model.VOUCHER_URL = config.Voucher.Link
	//GCS
	model.GCS_BUCKET = config.Gcs.Bucket
	model.GCS_PROJECT_ID = config.Gcs.ProjectID
	model.PUBLIC_URL = config.Gcs.PublicURL
	//OCRA
	model.OCRA_EVOUCHER_APPS_KEY = config.Ocra.AppsKey
	model.OCRA_URL = config.Ocra.Endpoint

	model.TOKENLIFE = config.Database.Redis.TokenLifetime

	negroni.NewLogger()
	r := setRoutes()
	m := negroni.New()
	m.Use(negroni.NewRecovery())
	m.Use(controller.LoggerMiddleware())
	m.Use(negroni.NewStatic(http.Dir(config.Server.PublicDirectory)))
	m.UseHandler(r)

	log.Printf("Server is listening on %q\n", config.Server.Host)
	http.ListenAndServe(config.Server.Host, m)
	//google app engine
	appengine.Main()
}

func printVersion() {
	fmt.Println(token)
	os.Exit(0)
}

func getUiRole() map[string][]string {
	m := make(map[string][]string)
	roles, err := model.GetAllUiFeatures()

	if err == nil {
		for _, value := range roles {
			m[value.Role] = append(m[value.Role], "/"+value.Category+"/"+value.Detail)
		}

	}
	//fmt.Print("Role ui : ")
	//fmt.Println(m)
	return m
}

func getApiRole() map[string][]string {
	m := make(map[string][]string)
	roles, err := model.GetAllApiFeatures()

	if err == nil {
		for _, value := range roles {
			m[value.Role] = append(m[value.Role], value.Category+"_"+value.Detail)
		}

	}
	//fmt.Print("Role api : ")
	//fmt.Println(m)
	return m
}

func getConfig() map[string]map[string]string {
	m := make(map[string]map[string]string)
	configs, err := model.GetAccountConfig()
	lastId := configs[0].AccountId
	mTemp := make(map[string]string)
	if err == nil {
		for i, value := range configs {
			if i != 0 && i+1 != len(configs) {
				lastId = configs[i+1].AccountId
				if value.AccountId != lastId {
					mTemp[value.ConfigDetail] = value.ConfigValue
					m[value.AccountId] = mTemp
					mTemp = make(map[string]string)
				}
			}
			mTemp[value.ConfigDetail] = value.ConfigValue
		}

	}
	m[lastId] = mTemp

	//fmt.Print("Config : ")
	//fmt.Println(m)
	return m
}
