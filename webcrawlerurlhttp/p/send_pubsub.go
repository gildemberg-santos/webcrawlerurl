package p

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
)

func SendPubSub(dataPubSub DataPubSub) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"))
	if err != nil {
		log.Println("pubsub.NewClient: ", err)
	}
	defer client.Close()

	t := client.Topic(os.Getenv("GOOGLE_TOPIC_NAME"))

	msg, err := json.Marshal(dataPubSub)
	if err != nil {
		log.Println("json.Marshal: ", err)
	}

	result := t.Publish(ctx, &pubsub.Message{
		Data: msg,
	})

	result.Ready()
	id, err := result.Get(ctx)
	if err != nil {
		log.Println("Error: result.Get -> ", err)
	}
	log.Println("Published a message; msg ID: ", id)
}
