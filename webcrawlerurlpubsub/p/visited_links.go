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
		UpdatedAt:  time.Now(),
		Validated:  true,
	}

	if !visitedLink.IsExist() {
		config.Logs("First Page")

		visitedLink.setDomain(url)
		visitedLink.init()
		visitedLink.saveOne()
		visitedLink.NormalizeLink()
		visitedLink.GetLink()
	}
}

func PendingPageLoop(loop int, company int32) {
	config.Logs("Pending Page Loop", loop)

	mongo := MongoDB{
		StringConnection: config.MongoStrConnec,
	}

	pendingPages, _ := mongo.FindAll(bson.M{"company": company, "status_link": "pending"}, bson.M{})
	var wg sync.WaitGroup
	for i, v := range pendingPages {
		wg.Add(1)
		go func(v bson.M) {
			visitedlink := VisitedLink{
				Company:    v["company"].(int32),
				Link:       v["link"].(string),
				Domain:     v["domain"].(string),
				StatusLink: v["status_link"].(string),
				Validated:  v["validated"].(bool),
				UpdatedAt:  v["updated_at"].(primitive.DateTime).Time(),
			}
			visitedlink.NormalizeLink()
			visitedlink.GetLink()
			defer wg.Done()
		}(v)
		if i >= config.LinksMax {
			break
		}
	}
	wg.Wait()

	retryPendingPages, _ := mongo.FindOne(bson.M{"company": company, "status_link": "pending"})
	if len(retryPendingPages) != 0 {
		loop += 1
		if loop <= config.LoopMax {
			PendingPageLoop(loop, company)
		}
		return
	}

}

func CleanDatabase(company int32) {
	config.Logs("Clean Database")

	mongo := MongoDB{
		StringConnection: config.MongoStrConnec,
	}

	databasePending, _ := mongo.FindAll(bson.M{"company": company, "status_link": "pending"}, bson.M{})
	databaseErrors, _ := mongo.FindAll(bson.M{"company": company, "status_link": "error"}, bson.M{})

	if len(databasePending) == 0 && len(databaseErrors) > 0 {
		mongo.DeleteAll(bson.M{"company": company, "status_link": "error"})
	}
}
