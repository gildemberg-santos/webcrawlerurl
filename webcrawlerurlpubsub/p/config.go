package p

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
)

type Config struct {
	LoopMax        int
	LinksMax       int
	MongoStrConnec string
	StatusLog      bool
	LimitCompany   int
}

type LimitCompanyError struct {
	Message    string
	Err        error
	StatusCode int
}

func (c *Config) Init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}

	config.SetLoopMax()
	config.SetLinksMax()
	config.SetMongoStrConnection()
	config.SetLimitCompany()
	config.SetLogs()
}

func (c *Config) SetLoopMax() {
	if os.Getenv("LOOP_MAX") != "" {
		intege, err := strconv.Atoi(os.Getenv("LOOP_MAX"))
		if err == nil && intege != 0 {
			c.LoopMax = intege
			return
		}
	}
	c.LoopMax = 1
}

func (c *Config) SetLinksMax() {
	if os.Getenv("LINKS_MAX") != "" {
		intege, err := strconv.Atoi(os.Getenv("LINKS_MAX"))
		if err == nil && intege != 0 {
			c.LinksMax = intege
			return
		}
	}
	c.LinksMax = 5
}

func (c *Config) SetMongoStrConnection() {
	if os.Getenv("MONGO_STR_CONNECTION") != "" {
		c.MongoStrConnec = os.Getenv("MONGO_STR_CONNECTION")
		return
	}
	c.MongoStrConnec = "mongodb://localhost:27017/"
}

func (c *Config) SetLogs() {
	if os.Getenv("LOG_SHOW") != "" {
		intege, err := strconv.Atoi(os.Getenv("LOG_SHOW"))
		if err == nil && intege != 0 {
			c.StatusLog = true
			return
		}
	}
	c.StatusLog = false
}

func (c *Config) SetLimitCompany() {
	if os.Getenv("LIMIT_COMPANY") != "" {
		intege, err := strconv.Atoi(os.Getenv("LIMIT_COMPANY"))
		if err == nil && intege != 0 {
			c.LimitCompany = intege
			return
		}
	}
	c.LimitCompany = 200
}

func (c *Config) IsLimitCompany(company int32) bool {
	c.SetMongoStrConnection()
	c.SetLimitCompany()

	mongo := MongoDB{
		StringConnection: config.MongoStrConnec,
	}

	pendingPages, _ := mongo.FindAll(bson.M{"company": company}, bson.M{})

	if len(pendingPages) >= c.LimitCompany {
		return true
	}
	return false
}

func (c *Config) Logs(msg ...interface{}) {
	if c.StatusLog {
		log.Println(msg...)
	}
}
