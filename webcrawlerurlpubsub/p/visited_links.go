package p

import (
	"log"
	"os"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FirstPage(company int32, url string) {
	visitedLink := VisitedLink{
		Company:    company,
		Link:       url,
		StatusLink: "pending",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Validated:  true,
	}

	if !visitedLink.IsExist() {
		log.Println("First Page")

		visitedLink.setDomain(url)
		visitedLink.init()
		visitedLink.GetLink()
	}
}

func PendingPageLoop(loop int32, company int32) {
	log.Println("Pending Page Loop", loop)

	mongo := MongoDB{
		StringConnection: os.Getenv("MONGO_STR_CONNECTION"),
	}

	pendingPages, _ := mongo.FindAll(bson.M{"company": company, "status_link": "pending"})
	var wg sync.WaitGroup
	for i, v := range pendingPages {
		wg.Add(1)
		visitedlink := VisitedLink{
			Company:    v["company"].(int32),
			Link:       v["link"].(string),
			Domain:     v["domain"].(string),
			StatusLink: v["status_link"].(string),
			Validated:  v["validated"].(bool),
			CreatedAt:  v["created_at"].(primitive.DateTime).Time(),
			UpdatedAt:  v["updated_at"].(primitive.DateTime).Time(),
		}

		go func(links VisitedLink) {
			links.GetLink()
			defer wg.Done()
		}(visitedlink)
		if i == 5 {
			break
		}
	}
	wg.Wait()
}

func CleanDatabase(company int32) {
	log.Println("Clean Database")

	mongo := MongoDB{
		StringConnection: os.Getenv("MONGO_STR_CONNECTION"),
	}

	databasePending, _ := mongo.FindAll(bson.M{"company": company, "status_link": "pending"})
	databaseErrors, _ := mongo.FindAll(bson.M{"company": company, "status_link": "error"})

	if len(databasePending) == 0 && len(databaseErrors) > 0 {
		mongo.DeleteAll(bson.M{"company": company, "status_link": "error"})
	}
}
