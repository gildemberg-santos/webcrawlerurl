package main

import (
	"net/http"
	"os"

	"github.com/gildemberg-santos/webcrawlerhttp/p"
)

func main() {
	os.Setenv("MONGO_STR_CONNECTION", "mongodb://localhost:27017") 
	os.Setenv("GOOGLE_CLOUD_PROJECT", "")
	os.Setenv("GOOGLE_TOPIC_NAME", "")
	http.HandleFunc("/", p.WebCrawlerUrlHttp)
	http.ListenAndServe(":8080", nil)
}
