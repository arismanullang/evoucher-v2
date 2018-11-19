package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/gilkor/evoucher/internal/model"
	"google.golang.org/api/option"

	"golang.org/x/net/context"
)

var (
	pscEvoucher *pubsub.Client
	pscJuno     *pubsub.Client
)

type PubSubMessage struct {
	Action string      `json:"action"`
	Data   interface{} `json:"data"`
}

type PubSubJunoAccount struct {
	Action string        `json:"action"`
	Data   PubSubAccount `json:"data"`
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

func AssignTenantPrivilegeVoucher() {
	var mu sync.Mutex

	topic := pscJuno.Topic("update-account")
	sub, err := createSubscriptionIfNotExists(pscJuno, topic, "update-account.privilege-voucher")
	if err != nil {
		log.Panic(err)
	}

	if err := sub.Receive(context.Background(), func(ctx context.Context, msg *pubsub.Message) {
		mu.Lock()
		defer mu.Unlock()

		var msgData PubSubJunoAccount
		if err := json.Unmarshal(msg.Data, &msgData); err != nil {
			log.Printf("Unable to process data: %v", err)
			msg.Ack()
			return
		}

		if msgData.Action != "create" {
			msg.Ack()
			return
		}

		data := msgData.Data
		gpr := model.GeneratePrivilegeRequest{
			CompanyID:  data.CompanyId,
			MemberID:   data.Id,
			MemberName: data.Name,
		}

		err := gpr.InsertPrivilegeVc()
		if err != nil {
			log.Println(err)
			msg.Nack()
			return
		}

		msg.Ack()
	}); err != nil {
		log.Fatal(err)
		return
	}
}

func InitPubSub() (err error) {
	ctx := context.Background()

	pscEvoucher, err = pubsub.NewClient(ctx, os.Getenv("GCLOUD_PROJECT"))
	if err != nil {
		return err
	}

	pscJuno, err = pubsub.NewClient(ctx, os.Getenv("JUNO_GCLOUD_PROJECT"), option.WithCredentialsFile(os.Getenv("JUNO_GCLOUD_CREDENTIALS")))
	if err != nil {
		return err
	}

	if _, err := createTopicIfNotExists(pscEvoucher, "update-vouchers"); err != nil {
		return err
	}

	return nil
}

func PublishDataTopic(data interface{}, topic string, action string) error {
	pubmsg := PubSubMessage{
		Action: action,
		Data:   data,
	}

	err := PublishTopic(pubmsg, topic)
	if err != nil {
		return err
	}

	return nil
}

func PublishTopic(pubmsg PubSubMessage, topicName string) error {
	topic := pscEvoucher.Topic(topicName)

	defer topic.Stop()

	msg, _ := json.Marshal(pubmsg)
	pubres := topic.Publish(context.Background(), &pubsub.Message{Data: msg})
	fmt.Println("Pubsub message : " + pubmsg.Action + string(msg))
	if _, err := pubres.Get(context.Background()); err != nil {
		return err
	}

	return nil
}

func createTopicIfNotExists(c *pubsub.Client, topic string) (*pubsub.Topic, error) {
	ctx := context.Background()

	t := c.Topic(topic)
	ok, err := t.Exists(ctx)
	if err != nil {
		return t, err
	}
	if ok {
		return t, nil
	}

	t, err = c.CreateTopic(ctx, topic)
	if err != nil {
		return t, err
	}

	return t, nil
}

func createSubscriptionIfNotExists(c *pubsub.Client, topic *pubsub.Topic, subscription string) (*pubsub.Subscription, error) {
	ctx := context.Background()

	s := c.Subscription(subscription)
	ok, err := s.Exists(ctx)
	if err != nil {
		return s, err
	}
	if ok {
		return s, nil
	}

	s, err = c.CreateSubscription(ctx, subscription, pubsub.SubscriptionConfig{
		Topic:       topic,
		AckDeadline: 20 * time.Second,
	})
	if err != nil {
		return s, err
	}
	return s, nil
}
