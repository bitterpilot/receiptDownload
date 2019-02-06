package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Config struct {
	Label   string `json:"Label"`
	Subject string `json:"Subject"`
	Sender  string `json:"Sender"`
}

func main() {
	config := loadConfiguration("./config.json")
	fmt.Println(config)
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
