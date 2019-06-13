package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/bitterpilot/receiptDownload/gmail"
	"github.com/bitterpilot/receiptDownload/link"
)

type config struct {
	Label            string `json:"Label"`
	Subject          string `json:"Subject"`
	Sender           string `json:"Sender"`
	SaveLoc          string `json:"SaveLoc"`
	ShortDescription string `json:"ShortDescription"`
	TimeZone         string `json:"TimeZone"`
}

type pdfInfo struct {
	Date string
	link.Link
}

func main() {
	config := loadConfiguration("config.json")
	query := fmt.Sprintf("from:%s subject:%s is:unread", config.Sender, config.Subject)

	var PDFList []pdfInfo
	list := gmail.ListEmails(config.Label, query)
	for _, msg := range list {
		msg = gmail.GetEmail(msg.Id)
		bodyDecoded, err := base64.URLEncoding.DecodeString(gmail.GetEmailBody(msg))
		if err != nil {
			log.Printf("Error decodong email: %v", err)
		}
		body := bytes.NewReader(bodyDecoded)
		links, err := link.Parse(body)
		if err != nil {
			log.Printf("Error finding links: %v", err)
		}

		// check that the link is to a PDF
		for _, link := range links {
			if strings.Contains(link.Text, ".pdf") {
				item := pdfInfo{
					Date: processInternalDate(msg.InternalDate, config.TimeZone),
					Link: link,
				}
				PDFList = append(PDFList, item)
			}
		}
	}

	// Write file
	for _, val := range PDFList {
		fileByte, receiptNum := getFile(val.Link.Href)
		// HACK: filename should probally have a space between the 2nd and 3rd %s. the lack of
		// 		 this space is a quick fick for the bug in the func getFile number variable
		fileName := fmt.Sprintf("%s %s%s.pdf", val.Date, config.ShortDescription, receiptNum)
		err := ioutil.WriteFile(fmt.Sprintf("%s/%s", config.SaveLoc, fileName), fileByte, 0644)
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

// code generated with postman(getpostman.com)
func getFile(url string) ([]byte, string) {
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("getFile:", err)
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	//get the receipt number from the filename
	// Make a Regex to say we only want numbers
	reg, err := regexp.Compile(`[^\d ]+`)
	if err != nil {
		log.Fatal(err)
	}
	filename := res.Header["Content-Disposition"]
	// BUG: this leaves a space in the number, the result is a file name with a
	// 		double space between ShortDescription and the number. When fixed
	// 		remove the hack in the write file for block in func main.
	number := reg.ReplaceAllString(strings.Join(filename, " "), "")
	return body, number
}

func processInternalDate(InternalDate int64, TimeZone string) string {
	tz, err := time.LoadLocation(TimeZone)
	if err != nil {
		log.Panicln(err)
	}
	dateUnix := time.Unix(0, InternalDate*int64(time.Millisecond)).In(tz)
	return dateUnix.Format("20060102")
}
