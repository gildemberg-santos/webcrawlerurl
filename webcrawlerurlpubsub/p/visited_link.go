package p

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"go.mongodb.org/mongo-driver/bson"
)

type VisitedLink struct {
	Company    int32     `bson:"company"`
	Link       string    `bson:"link"`
	StatusLink string    `bson:"status_link"`
	Validated  bool      `bson:"validated"`
	Domain     string    `bson:"domain"`
	UpdatedAt  time.Time `bson:"updated_at"`
}

func (v *VisitedLink) init() {
	if v.Company == 0 {
		config.Logs("Company is required")
		v.Validated = false
	}

	if v.Link == "" {
		config.Logs("Link is required")
		v.Validated = false
	}

	if v.StatusLink == "" {
		config.Logs("StatusLink is required")
		v.Validated = false
	}

	v.NormalizeLink()
}

func (v *VisitedLink) GetLink() {
	v.StatusLink = "visited"
	v.Validated = true
	v.saveOne()

	if !v.Validated {
		config.Logs("Error: Invalid link", v.Link)
		return
	}

	resp, err := http.Get(v.Link)
	if err != nil {
		config.Logs("Error: On get link", v.Link)
		return
	}
	defer resp.Body.Close()

	var links = v.extractLinks(resp.Body)
	v.saveMany(links)
}

func (v *VisitedLink) extractLinks(node io.Reader) []string {
	if !v.Validated {
		return []string{}
	}

	doc, err := goquery.NewDocumentFromReader(node)
	if err != nil {
		config.Logs("Error: On parse link", v.Link)
		return []string{}
	}

	links := []string{}
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, success := s.Attr("href")

		if success && href != "" {
			var visitedLink = VisitedLink{
				Link:      href,
				Domain:    v.Domain,
				Validated: true,
			}
			visitedLink.NormalizeLink()
			if visitedLink.Validated {
				links = append(links, visitedLink.Link)
			}
		}
	})
	return links
}

func (v *VisitedLink) NormalizeLink() {
	if !v.Validated {
		return
	}

	link, err := url.Parse(v.Link)

	if err != nil {
		v.Validated = false
		return
	}

	if link.Host == "" {
		link.Scheme = "https"
		link.Host = v.Domain
	}

	if link.Scheme == "" {
		link.Scheme = "https"
	}

	if link.Host != v.Domain {
		sub_host := strings.LastIndex(link.Host, v.Domain)
		if sub_host == -1 {
			v.Validated = false
			return
		}
	}

	if link.Path != "" {
		validationExtension := func(path string, valid string) bool {
			return strings.HasSuffix(path, valid)
		}

		validationWord := func(path string, valid string) bool {
			return strings.LastIndex(path, valid) != -1
		}

		for _, extension := range []string{".pdf", ".jpg", ".gif", ".png"} {
			if validationExtension(link.Path, extension) {
				v.Validated = false
				return
			}
		}

		for _, word := range []string{"mailto:", "tel:", "javascript:", "window."} {
			if validationWord(link.Path, word) {
				v.Validated = false
				return
			}
		}

		if validationWord(link.Path, v.Domain) {
			v.Validated = false
			return
		}
	}

	link.Fragment = ""
	link.RawQuery = ""

	v.Link = link.String()
	v.Validated = true
}

func (v *VisitedLink) saveOne() {
	if !v.Validated {
		return
	}

	mongo := MongoDB{
		StringConnection: config.MongoStrConnec,
	}

	if mongo.StringConnection == "" {
		config.Logs("StringConnection is empty")
		return
	}

	linksDB, err := mongo.FindOne(bson.M{"link": v.Link, "company": v.Company})

	if err == nil {
		if linksDB["status_link"] == "visited" {
			return
		}

		if linksDB["company"] == v.Company && linksDB["link"] == v.Link && linksDB["status_link"] == v.StatusLink {
			return
		}
	}

	if v.StatusLink == "visited" {
		v.UpdatedAt = time.Now()
		mongo.UpsertOne(v, bson.M{"link": v.Link, "company": v.Company})
	} else if v.StatusLink == "pending" {
		v.UpdatedAt = time.Now()
		mongo.InsertOne(v)
	}

	config.Logs("Saved", "status", v.StatusLink, "company", v.Company, "link", v.Link)
}

func (v *VisitedLink) saveMany(links []string) {
	if !v.Validated {
		return
	}

	mongo := MongoDB{
		StringConnection: config.MongoStrConnec,
	}

	if mongo.StringConnection == "" {
		config.Logs("StringConnection is empty")
		return
	}

	if v.StatusLink == "visited" {
		visitedLinks := bson.A{}
		for _, link := range links {
			visitedLink := bson.M{
				"link":        link,
				"company":     v.Company,
				"domain":      v.Domain,
				"status_link": "pending",
				"validated":   true,
				"created_at":  time.Now(),
				"updated_at":  time.Now(),
			}
			visitedLinks = append(visitedLinks, visitedLink)
		}

		mongo.InsertMany(visitedLinks)
		config.Logs("Saved", "status", v.StatusLink, "company", v.Company, "Inserted", len(visitedLinks))
	}
}

func (v *VisitedLink) IsExist() bool {
	mongo := MongoDB{
		StringConnection: config.MongoStrConnec,
	}

	if mongo.StringConnection == "" {
		config.Logs("StringConnection is empty")
		return false
	}

	linksDB, err := mongo.FindOne(bson.M{"link": v.Link, "company": v.Company})

	if err != nil {
		return false
	}

	if linksDB["status_link"] == "visited" {
		return true
	}

	return false
}

func (v *VisitedLink) setDomain(domain string) {
	if !v.Validated {
		return
	}

	link, err := url.Parse(domain)

	if err != nil {
		v.Validated = false
		return
	}

	v.Domain = link.Host
}
