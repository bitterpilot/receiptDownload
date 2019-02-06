package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/bitterpilot/receiptDownload/gmail"
)

type Config struct {
	Label   string `json:"Label"`
	Subject string `json:"Subject"`
	Sender  string `json:"Sender"`
	SaveLoc string `json:"SaveLoc"`

	//temporary example items
	Example struct {
		EmailID string `json:"emailID"`
	} `json:"example"`
}

func main() {
	config := loadConfiguration("config.json")
	query := fmt.Sprintf("from:%s subject:%s", config.Sender, config.Subject)
	list := gmail.ListEmails(config.Label, query)

	if len(list) > 3 {
		//GetEmails
	} else {
		for _, val := range list {
			fmt.Println(val.Id)
			//GetEmail
		}
	}
	msg := gmail.GetEmail(config.Example.EmailID)
}

func loadConfiguration(file string) Config {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}

	jsonParser := json.NewDecoder(configFile)
	if err := jsonParser.Decode(&config); err != nil {
		log.Fatalf("json error: %s", err)
	}
	return config
}
