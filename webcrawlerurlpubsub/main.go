package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/gildemberg-santos/webcrawlerurlpubsub/p"
)

func main() {
	os.Setenv("LOG_SHOW", "0")
	var cxt = context.Context(context.Background())
	var datalinks = p.DataLinks{
		Company: -1,
		Link:    "https://www.iteva.org.br",
	}
	var m = p.PubSubMessage{}
	bytesjson, _ := json.Marshal(datalinks)
	m.Data = bytesjson

	p.WebCrawlerUrlPubSub(cxt, m)

}
