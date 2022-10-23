package p

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DataResult struct {
	Company int32         `json:"company"`
	Link    string        `json:"link"`
	Result  []primitive.M `json:"result"`
	Total   int32         `json:"total"`
}

type DataPubSub struct {
	Company int32  `json:"company"`
	Link    string `json:"link"`
}

func WebCrawlerUrlHttp(w http.ResponseWriter, r *http.Request) {
	var d = DataResult{}

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		switch err {
		case io.EOF:
			http.Error(w, http.StatusText(http.StatusBadRequest)+": Nenhuma informação recebida", http.StatusBadRequest)
			return
		default:
			log.Printf("json.NewDecoder: %v", err)
			http.Error(w, http.StatusText(http.StatusBadRequest)+": json.NewDecoder", http.StatusBadRequest)
			return
		}
	}

	if d.Link == "" || d.Company == 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest)+": Nenhuma informação encontrada", http.StatusBadRequest)
		return
	}

	mongo := MongoDB{
		StringConnection: os.Getenv("MONGO_STR_CONNECTION"),
	}
	d.Result, _ = mongo.FindAll(bson.M{"company": d.Company})
	if d.Result == nil {
		d.Result = []primitive.M{}
	}
	d.Total = int32(len(d.Result))

	if os.Getenv("GOOGLE_CLOUD_PROJECT") != "" && os.Getenv("GOOGLE_TOPIC_NAME") != ""{
		SendPubSub(DataPubSub{Company: d.Company, Link: d.Link})
	}

	json.NewEncoder(w).Encode(d)
}
