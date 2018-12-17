package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/ruizu/render"
	"github.com/urfave/negroni"

	"github.com/gilkor/evoucher/internal/controller"
	"github.com/gilkor/evoucher/internal/model"
	"github.com/gilkor/evoucher/lib/server"
)

func main() {

	if err := model.ConnectDB(os.Getenv("DB")); err != nil {
		log.Fatal(err)
	}

	render.SetPath(os.Getenv("TEMPLATE_DIR"))

	model.UiFeatures = getUiRole()
	model.ApiFeatures = getApiRole()
	model.Config = getConfig()
	model.Domain = os.Getenv("MAILGUN_DOMAIN")
	model.ApiKey = os.Getenv("MAILGUN_KEY")
	model.PublicApiKey = os.Getenv("MAILGUN_PUBLIC_KEY")
	model.RootTemplate = os.Getenv("MAILGUN_TEMPLATE_DIR")
	model.Email = os.Getenv("MAILGUN_FROM")
	model.RootURL = os.Getenv("MAILGUN_ROOT_URL")

	model.GetProgramTypes()

	//voucher config
	model.VOUCHER_URL = os.Getenv("VOUCHER_LINK")
	//GCS
	model.GCS_BUCKET = os.Getenv("GCS_BUCKET")
	model.GCS_PROJECT_ID = os.Getenv("GCLOUD_PROJECT")
	model.PUBLIC_URL = os.Getenv("GCS_PUBLIC_URL")

	tokenTTL, _ := strconv.Atoi(os.Getenv("TOKEN_TTL"))
	model.TOKENLIFE = tokenTTL

	m := negroni.New()
	m.Use(controller.LoggerMiddleware())
	m.Use(negroni.NewStatic(http.Dir(os.Getenv("PUBLIC_DIR"))))
	m.UseHandler(router)

	err := controller.InitPubSub()
	if err != nil {
		log.Fatal(err)
	}

	go controller.AssignTenantPrivilegeVoucher()

	log.Fatal(server.ListenAndServe(m))

}

func getUiRole() map[string][]string {
	m := make(map[string][]string)
	roles, err := model.GetAllUiFeatures()

	if err == nil {
		for _, value := range roles {
			m[value.Role] = append(m[value.Role], "/"+value.Category+"/"+value.Detail)
		}

	}
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

	return m
}
