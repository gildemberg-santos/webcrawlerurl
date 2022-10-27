package p

import (
	"log"
	"os"
	"strconv"
)

func LoopMax() int {
	if os.Getenv("LOOP_MAX") != "" {
		intege, err := strconv.Atoi(os.Getenv("LOOP_MAX"))
		if err == nil && intege == 1 {
			return intege
		}
	}
	return 1
}

func LinksMax() int {
	if os.Getenv("LINKS_MAX") != "" {
		intege, err := strconv.Atoi(os.Getenv("LINKS_MAX"))
		if err == nil && intege == 1 {
			return intege
		}
	}
	return 5
}

func MongoStrConnection() string {
	if os.Getenv("MONGO_STR_CONNECTION") != "" {
		return os.Getenv("MONGO_STR_CONNECTION")
	}
	return "mongodb://localhost:27017/"
}

func Logs(msg ...interface{}) {
	if os.Getenv("LOG_SHOW") != "" {
		intege, err := strconv.Atoi(os.Getenv("LOG_SHOW"))
		if err == nil && intege == 1 {
			log.Println(msg...)
		}
	}
}
