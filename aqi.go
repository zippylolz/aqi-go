package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

/**
* AQI Category and Value
* 	Good: 0-50 green
* 	Moderate: 51-100 yellow
* 	Unhealthy for Sensitive Groups: 101-150 light orange
* 	Unhealthy: 151-200 dark orange
* 	Very Unhealthy: 201-300 light red
* 	Hazardous: 301+ dark red
 */

type Rss struct {
	XMLName xml.Name `xml:"rss"`
	Aqi     Aqi      `xml:"channel"`
}

type Aqi struct {
	XMLName     xml.Name `xml:"channel"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
	Language    string   `xml:"language"`
	WebMaster   string   `xml:"webMaster"`
	PublishDate string   `xml:"pubDate"`
	Item        AqiItem  `xml:"item"`
}

type AqiItem struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
}

func main() {

	feedUrlPtr := flag.String(
		"url",
		"https://feeds.airnowapi.org/rss/realtime/24.xml",
		"Go to https://feeds.airnowapi.org/ to find the Current Air Quality rss feed.",
	)

	flag.Parse()

	// Download xml
	resp, err := http.Get(*feedUrlPtr)

	if err != nil {
		fmt.Println(err)
	}

	// read xmlFile as a byte array.
	byteValue, _ := io.ReadAll(resp.Body)

	// initialize Rss struct
	var rss Rss

	xml.Unmarshal(byteValue, &rss)

	// parse html to extract particle pollution aqi number
	particalRe, err := regexp.Compile(`(?s)\d{1,3}\s?AQI\s?-\s?Particle Pollution`)

	if err != nil {
		fmt.Println(err)
	}

	var sb strings.Builder
	sb.WriteString(rss.Aqi.Title + "\n")

	var levelRe = regexp.MustCompile(`(?m)[<>\W\da-zA-Z\t\n:;\" ]*<br\s/><br\s/>\s*<div>\s+(?P<Level>[\w\s]*)\s-[<>\W\da-zA-Z\t\n ]*`)

	matches := levelRe.FindStringSubmatch(rss.Aqi.Item.Description)

	if len(matches) > 0 {
		levelIndex := levelRe.SubexpIndex("Level")
		sb.WriteString(matches[levelIndex])
	} else {
		fmt.Println("No level found.")
	}

	// particle pollution
	if len(particalRe.FindStringIndex(rss.Aqi.Item.Description)) > 0 {
		var particleValue = strings.Split(particalRe.FindString(rss.Aqi.Item.Description), " ")[0]
		if sb.Len() > 0 {
			sb.WriteString("\t")
		}
		sb.WriteString(particleValue + " AQI" + "\tPM2.5")
	}

	fmt.Println(sb.String())
}
