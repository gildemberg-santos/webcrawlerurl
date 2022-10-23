package p

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VisitedLink struct {
	Company    int32     `bson:"company"`
	Link       string    `bson:"link"`
	StatusLink string    `bson:"status_link"`
	StatusCode int32     `bson:"status_code"`
	Validated  bool      `bson:"validated"`
	Domain     string    `bson:"domain"`
	CreatedAt  time.Time `bson:"created_at"`
	UpdatedAt  time.Time `bson:"updated_at"`
}

func (v *VisitedLink) init() {
	if v.Company == 0 {
		log.Println("Company is required")
		v.Validated = false
	}

	if v.Link == "" {
		log.Println("Link is required")
		v.Validated = false
	}

	if v.StatusLink == "" {
		log.Println("StatusLink is required")
		v.Validated = false
	}

	v.normalizeLink()

	if v.Validated {
		v.save()
	}
}

func (v *VisitedLink) GetLink() {
	v.init()

	if !v.Validated {
		log.Println("Error: Invalid link", v.Link)
		return
	}

	resp, err := http.Get(v.Link)
	if err != nil {
		log.Println("Error: On get link", v.Link)
		v.Validated = true
		v.StatusLink = "error"
		v.save()
		return
	}
	defer resp.Body.Close()

	v.StatusCode = int32(resp.StatusCode)

	if resp.StatusCode == http.StatusNotFound {
		v.Validated = true
		v.StatusLink = "error"
		v.save()
		return
	}

	v.StatusLink = "visited"
	v.save()

	v.extractLinks(resp.Body)
}

func (v *VisitedLink) extractLinks(node io.Reader) {
	if !v.Validated {
		return
	}

	doc, err := goquery.NewDocumentFromReader(node)
	if err != nil {
		log.Println("Error: On parse link", v.Link)
		v.Validated = true
		v.StatusLink = "error"
		v.save()
		return
	}

	links := []string{}
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, success := s.Attr("href")

		if success && href != "" {
			visitedLink := VisitedLink{Link: href}
			visitedLink.normalizeLink()
			links = append(links, visitedLink.Link)
		}
	})

	var wg sync.WaitGroup
	for _, link := range links {
		wg.Add(1)
		go func(link string) {
			visitedLink := VisitedLink{
				Company:    v.Company,
				Link:       link,
				StatusLink: "pending",
				Domain:     v.Domain,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
				Validated:  true,
			}
			visitedLink.normalizeLink()
			if visitedLink.Validated {
				visitedLink.init()
			}
			defer wg.Done()
		}(link)
	}
	wg.Wait()
}

func (v *VisitedLink) normalizeLink() {
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

func (v *VisitedLink) save() {
	if !v.Validated {
		return
	}

	mongo := MongoDB{
		StringConnection: os.Getenv("MONGO_STR_CONNECTION"),
	}

	if mongo.StringConnection == "" {
		log.Println("StringConnection is empty")
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

		v.CreatedAt = linksDB["created_at"].(primitive.DateTime).Time()
		v.UpdatedAt = time.Now()
	}

	mongo.Upsert(v, bson.M{"link": v.Link, "company": v.Company})
	log.Println("Saved", "status", v.StatusLink, "company", v.Company, "link", v.Link)
}

func (v *VisitedLink) IsExist() bool {
	mongo := MongoDB{
		StringConnection: os.Getenv("MONGO_STR_CONNECTION"),
	}

	if mongo.StringConnection == "" {
		log.Println("StringConnection is empty")
		return false
	}

	linksDB, err := mongo.FindOne(bson.M{"link": v.Link, "company": v.Company})

	if err != nil {
		return false
	}

	if linksDB["status_link"] == "visited" || linksDB["status_link"] == "error" {
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

func linksUnique(links []string) []string {
	linksTemp := []string{}
	validatedLink := map[string]bool{}

	for v := range links {
		linkTemp := VisitedLink{
			Link: links[v],
		}
		linkTemp.normalizeLink()
		validatedLink[linkTemp.Link] = true
	}

	for key := range validatedLink {
		linkTemp := VisitedLink{
			Link: key,
		}
		linkTemp.normalizeLink()
		linksTemp = append(linksTemp, linkTemp.Link)
	}

	return linksTemp
}

func linksExist(links []string, link string) bool {
	encountered := map[string]bool{}

	for v := range links {
		encountered[links[v]] = true
	}

	return encountered[link]
}
