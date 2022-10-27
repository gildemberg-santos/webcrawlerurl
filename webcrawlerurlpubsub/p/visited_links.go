package p

import (
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
		Logs("First Page")

		visitedLink.setDomain(url)
		visitedLink.init()
		visitedLink.saveOne()
		visitedLink.GetLink()
	}
}

func PendingPageLoop(loop int, company int32) {
	Logs("Pending Page Loop", loop)

	mongo := MongoDB{
		StringConnection: MongoStrConnection(),
	}

	pendingPages, _ := mongo.FindAll(bson.M{"company": company, "status_link": "pending"}, bson.M{})
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
		if i >= LinksMax() {
			break
		}
	}
	wg.Wait()

	retryPendingPages, _ := mongo.FindOne(bson.M{"company": company, "status_link": "pending"})
	if len(retryPendingPages) != 0 {
		loop += 1
		if loop > LoopMax() {
			return
		}
		time.Sleep(1 * time.Second)
		PendingPageLoop(loop, company)
	}

}

func CleanDatabase(company int32) {
	Logs("Clean Database")

	mongo := MongoDB{
		StringConnection: MongoStrConnection(),
	}

	databasePending, _ := mongo.FindAll(bson.M{"company": company, "status_link": "pending"}, bson.M{})
	databaseErrors, _ := mongo.FindAll(bson.M{"company": company, "status_link": "error"}, bson.M{})

	if len(databasePending) == 0 && len(databaseErrors) > 0 {
		mongo.DeleteAll(bson.M{"company": company, "status_link": "error"})
	}
}
