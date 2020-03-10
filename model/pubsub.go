package model

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"sync"

	"cloud.google.com/go/pubsub"

	"github.com/gilkor/athena/lib/x"
)

var (
	psc *pubsub.Client
)

func init() {
	var err error
	// connect to Cloud Pub/Sub
	ctx := context.Background()
	psc, err = pubsub.NewClient(ctx, os.Getenv("GCLOUD_PUBSUB_PROJECT"))
	if err != nil {
		log.Fatal(err)
	}

	x.FatalOnError(createPubSubSubscribtion("update-account", "update-account.evoucher-v2"))
}

func createPubSubSubscribtion(topic string, subscriptions ...string) error {
	ctx := context.Background()

	for _, s := range subscriptions {
		exists, err := psc.Subscription(s).Exists(ctx)
		if !exists && err == nil {
			_, err := psc.CreateSubscription(ctx, s, pubsub.SubscriptionConfig{
				Topic:       psc.Topic(topic),
				AckDeadline: 0,
			})
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
	}

	return nil
}

func StartSubscriber() {
	go startAccountSubscriber()
}

type PubSubAccount struct {
	Action  string  `json:"action"`
	Account Account `json:"data"`
}

func startAccountSubscriber() {
	var mu sync.Mutex
	sub := psc.Subscription("update-account.evoucher-v2")
	if err := sub.Receive(context.Background(), func(ctx context.Context, msg *pubsub.Message) {
		mu.Lock()
		defer mu.Unlock()

		var data PubSubAccount
		if err := json.Unmarshal(msg.Data, &data); err != nil {
			log.Printf("Unable to process data: %v", err)
			msg.Ack()
			return
		}

		ac := data.Account

		switch data.Action {
		case "create":
			if _, err := ac.Insert(); err != nil {
				log.Printf("Error when insert account with id: %v", ac.ID)
				return
			}
			log.Printf("Innsert account with id: %v", ac.ID)

		case "delete":
			if err := ac.Update(); err != nil {
				log.Printf("Error when delete account with id: %v", ac.ID)
			}
			log.Printf("Delete account with id: %v", ac.ID)

		default:
			if err := ac.Update(); err != nil {
				if _, err := ac.Insert(); err != nil {
					log.Printf("Error when update account with id: %v", ac.ID)
					return
				}
			}
			log.Printf("Update account with id: %v", ac.ID)
		}

		msg.Ack()
	}); err != nil {
		log.Fatal(err)
		return
	}
}
