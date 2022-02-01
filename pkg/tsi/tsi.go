package tsi

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type TSI struct {
	FeedUrl string
}

func NewTSI() *TSI {
	// TODO: implement getting the feed url from the config
	return &TSI{FeedUrl: "https://www.tsi.lv/ru/feed/"}
}

type News struct {
	Titles []string `xml:"channel>item>title"`
	Links  []string `xml:"channel>item>link"`
	Dates  []string `xml:"channel>item>pubDate"`
}

func (t TSI) GetNews() (string, error) {
	news, err := t.getNews()
	if err != nil {
		return "", err
	}
	return t.parseNews(news), nil
}

func (t TSI) getNews() (News, error) {
	news, err := t.parseFeed(t.FeedUrl)
	if err != nil {
		return News{}, err
	}
	return news, nil
}

func (t TSI) parseFeed(feedUrl string) (News, error) {
	resp, err := http.Get(feedUrl)
	if err != nil {
		return News{}, err
	}

	defer resp.Body.Close()

	byteValue, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return News{}, err
	}

	news := News{}
	err = xml.Unmarshal(byteValue, &news)
	if err != nil {
		log.Fatal(err)
		return News{}, err
	}

	log.Println(news.Titles[0])

	return news, nil
}

func (t TSI) parseNews(news News) string {
	var newsString string
	for i := 0; i < 5; i++ {
		newsString += t.parseDate(news.Dates[i]) + "\n"
		newsString += news.Titles[i] + "\n"
		newsString += news.Links[i] + "\n\n"
	}
	return newsString
}

func (t TSI) parseDate(date string) string {
	parse, err := time.Parse(time.RFC1123Z, date)
	if err != nil {
		return ""
	}
	return parse.Format("02.01.2006")
}
