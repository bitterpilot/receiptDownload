package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/html"

	"github.com/bitterpilot/receiptDownload/gmail"
)

type config struct {
	Label   string `json:"Label"`
	Subject string `json:"Subject"`
	Sender  string `json:"Sender"`
	SaveLoc string `json:"SaveLoc"`
}

func main() {
	config := loadConfiguration("config.json")
	query := fmt.Sprintf("from:%s subject:%s", config.Sender, config.Subject)
	list := gmail.ListEmails(config.Label, query)

	var linklist []string
	for _, msg := range list {
		//https://www.thepolyglotdeveloper.com/2017/05/concurrent-golang-applications-goroutines-channels/
		link, _ := getLink(gmail.GetEmailBody(gmail.GetEmail(msg.Id)))
		linklist = append(linklist, link)
	}

	destination := config.SaveLoc
	for _, val := range linklist {
		fileByte := getFile(val)
		// generate a unique name using the url.
		// the url includes .pdf as well so remove that to begin with
		// the future will look something like
		// Sprintf("%s%s", dateFromEmail, ".pdf")
		fileName := fmt.Sprintf("%s%s", val[len(val)-10:len(val)-4], ".pdf")
		err := ioutil.WriteFile(fmt.Sprintf("%s/%s", destination, fileName), fileByte, 0644)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func loadConfiguration(file string) config {
	var conf config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}

	jsonParser := json.NewDecoder(configFile)
	if err := jsonParser.Decode(&conf); err != nil {
		log.Fatalf("json error: %s", err)
	}
	return conf
}

func getLink(body string) (link string, err error) {
	decode, err := base64.URLEncoding.DecodeString(body)
	if err != nil {
		fmt.Println(err)
	}
	doc, err := html.Parse(bytes.NewReader(decode))
	if err != nil {
		fmt.Println(err)
	}

	var b *html.Node
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			b = n
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	for _, attr := range b.Attr {
		if attr.Key == "href" {
			return attr.Val, nil
		}
	}
	return "", errors.New("no link found")
}

// thank you postman(getpostman.com)
func getFile(url string) []byte {
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("cache-control", "no-cache")

	// filename is in here
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("getFile:", err)
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	return body
}
