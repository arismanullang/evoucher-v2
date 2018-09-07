package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/ruizu/render"
	"github.com/urfave/negroni"

	"github.com/gilkor/evoucher/internal/controller"
	"github.com/gilkor/evoucher/internal/model"
	"github.com/gilkor/evoucher/lib/server"
)

var psc *pubsub.Client

func init() {
	var err error
	ctx := context.Background()
	psc, err = pubsub.NewClient(ctx, os.Getenv("GCLOUD_PROJECT"))
	if err != nil {
		log.Fatal(err)
	}
}

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

	go assignTenantPrivilegeVoucher()

	log.Fatal(server.ListenAndServe(m))

}

type PubSubAccount struct {
	Id                 string      `json:"id"`
	Name               string      `json:"name"`
	CompanyId          string      `json:"company_id"`
	Gender             string      `json:"gender"`
	BirthDate          string      `json:"birthdate"`
	BrithPlace         string      `json:"birthplace"`
	MaritalStatus      string      `json:"marital_status"`
	IdentityNo         string      `json:"identity_no"`
	IdentityType       string      `json:"identity_type"`
	IdentityIssuedDate string      `json:"identity_issue_date"`
	OccupationId       string      `json:"occupation_id"`
	ReligionId         string      `json:"religion_id"`
	Email              string      `json:"email"`
	MobileCallingCode  string      `json:"mobile_calling_code"`
	MobileNo           string      `json:"mobile_no"`
	Address            string      `json:"address"`
	CountryCode        string      `json:"country_code"`
	StateId            string      `json:"state_id"`
	CityId             string      `json:"city_id"`
	DistrictId         string      `json:"district_id"`
	VillageId          string      `json:"village_id"`
	ZipCode            string      `json:"zip_code"`
	State              string      `json:"state"`
	CreatedBy          string      `json:"created_by"`
	CreatedAt          time.Time   `json:"created_at"`
	UpdatedBy          string      `json:"updated_by"`
	UpdatedAt          time.Time   `json:"updated_at"`
	DeletedBy          string      `json:"deleted_by"`
	DeletedAt          interface{} `json:"deleted_at"`
}

func assignTenantPrivilegeVoucher() {
	var mu sync.Mutex
	sub := psc.Subscription("update-account")
	if err := sub.Receive(context.Background(), func(ctx context.Context, msg *pubsub.Message) {
		mu.Lock()
		defer mu.Unlock()

		var data PubSubAccount
		if err := json.Unmarshal(msg.Data, &data); err != nil {
			log.Printf("Unable to process data: %v", err)
			msg.Ack()
			return
		}

		gpr := model.GeneratePrivilegeRequest{
			CompanyID:  data.CompanyId,
			MemberID:   data.Id,
			MemberName: data.Name,
		}

		err := gpr.InsertPrivilegeVc()
		if err != nil {
			msg.Ack()
		}

		msg.Ack()
	}); err != nil {
		log.Fatal(err)
		return
	}
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
