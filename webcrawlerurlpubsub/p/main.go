package p

import (
	"context"
	"encoding/json"
	"log"
)

var config = Config{}

type PubSubMessage struct {
	Data []byte `json:"data"`
}

type DataLinks struct {
	Company int32  `json:"company"`
	Link    string `json:"link"`
}

func WebCrawlerUrlPubSub(ctx context.Context, m PubSubMessage) error {
	config.Init()
	var datalinks DataLinks
	json.Unmarshal(m.Data, &datalinks)
	call(datalinks.Company, datalinks.Link)
	return nil
}

func call(company int32, link string) {
	log.Println("Starting webcrawlerurl company", company, "link", link)
	FirstPage(company, link)
	PendingPageLoop(1, company)
	CleanDatabase(company)
	log.Println("Done webcrawlerurl company", company, "link", link)
}
