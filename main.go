package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/html"
)

func main() {
	fmt.Println("Starting webcrawler")
	done := make(chan bool)
	go VisitLink(-1, "https://blog.rocketseat.com.br/", done, 1)
	<-done
	fmt.Println("Finished webcrawler")
}

// *****************************************************************************
// Extrair e visitar links para salvar no banco de dados
// struct: VisitedLink,
// func: VisitLink, extractLinks
// *****************************************************************************

type VisitedLink struct {
	Company     int       `json:"company" bson:"company"`
	Website     string    `json:"website"`
	Link        string    `json:"link"`
	VisitedDate time.Time `json:"visited_date"`
}

func VisitLink(company int, link string, done chan bool, level int) {
	if company == 0 {
		fmt.Println("Error: No companies reported.")
		done <- true
		return
	}

	link_uri, err := url.Parse(link)
	if err != nil || link_uri.Scheme == "" {
		fmt.Println("Error: Invalid link", link)
		return
	}

	resp, err := http.Get(link_uri.String())
	if err != nil {
		fmt.Println("Error: On get link", link)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: Status code ", resp.StatusCode, " on link", link)
		return
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Println("Error: On parse link", link)
		return
	}

	extractLinks(doc, link_uri.Host, company, done)

	if level == 1 {
		done <- true
	}
}

func extractLinks(node *html.Node, host string, company int, done chan bool) {
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, attr := range node.Attr {
			if attr.Key != "href" {
				continue
			}

			link, err := url.Parse(attr.Val)

			if err != nil {
				continue
			}

			if link.Host == "" {
				link.Scheme = "https"
				link.Host = host
			}

			if link.Host != host || link.Scheme == "" {
				continue
			}

			if link.Path != "" {
				extension := strings.LastIndex(link.Path, ".")
				mailto := strings.LastIndex(link.Path, "mailto:")
				tel := strings.LastIndex(link.Path, "tel:")
				javascript := strings.LastIndex(link.Path, "javascript:")
				window := strings.LastIndex(link.Path, "window.")

				if extension != -1 || mailto != -1 || tel != -1 || javascript != -1 || window != -1 {
					link.Path = ""
				}
			}

			link.Fragment = ""
			link.RawQuery = ""

			if CheckVisitedLink(link.String()) {
				fmt.Println("Link already visited", link.String())
				continue
			}

			visitedlink := VisitedLink{
				Company:     company,
				Website:     link.Host,
				Link:        link.String(),
				VisitedDate: time.Now(),
			}

			go VisitLink(company, link.String(), done, 0)
			Insert("links", visitedlink)

			fmt.Println("Inserted link", link.String())

		}
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		extractLinks(c, host, company, done)
	}
}

// *****************************************************************************
// MongoDB
// func: getConnection, Insert, CheckVisitedLink
// *****************************************************************************

func getConnection() (client *mongo.Client, ctx context.Context) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ = context.WithTimeout(context.Background(), 15*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return
}

func Insert(collection string, data interface{}) error {
	client, ctx := getConnection()
	defer client.Disconnect(ctx)

	c := client.Database("crawler").Collection(collection)
	_, err := c.InsertOne(context.Background(), data)

	return err
}

func CheckVisitedLink(link string) bool {
	client, ctx := getConnection()
	defer client.Disconnect(ctx)

	c := client.Database("crawler").Collection("links")
	opts := options.Count().SetLimit(1)
	n, err := c.CountDocuments(context.TODO(), bson.D{{"link", link}}, opts)
	if err != nil {
		fmt.Println("Error MongoDB", err)
		return true
	}

	return n > 0
}
