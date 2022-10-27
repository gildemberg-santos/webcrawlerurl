package p

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	LoopMax        int
	LinksMax       int
	MongoStrConnec string
	StatusLog      bool
}

func (c *Config) Init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	config.SetLoopMax()
	config.SetLinksMax()
	config.SetMongoStrConnection()
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

func (c *Config) Logs(msg ...interface{}) {
	if c.StatusLog {
		log.Println(msg...)
	}
}
