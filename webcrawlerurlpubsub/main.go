package main

import (
	"context"
	"encoding/json"

	"github.com/gildemberg-santos/webcrawlerurlpubsub/p"
)

func main() {
	var cxt = context.Context(context.Background())
	var datalinks = p.DataLinks{
		Company: 4,
		Link:    "https://www.iteva.org.br/",
	}
	var m = p.PubSubMessage{}
	bytesjson, _ := json.Marshal(datalinks)
	m.Data = bytesjson

	p.WebCrawlerUrlPubSub(cxt, m)

}
