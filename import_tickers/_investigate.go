package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

type Dst struct {
	Symbol       string
	Title        string
	SimplyWallSt string
}

func main() {
	Run()
}

func Run() bool {

	result := make(map[string]*Dst)

	c := colly.NewCollector()
	c.SetRequestTimeout(4 * time.Minute)
	c.DisableCookies()
	c.IgnoreRobotsTxt = true
	c.AllowURLRevisit = true

	industry := "" // "/banks"
	page := 1
	count := 10

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if isTicker(e.Attr("href")) {
			log.Print(e.ChildText("b"))
			result[e.ChildText("b")] = &Dst{
				Symbol:       e.ChildText("b"),
				Title:        e.ChildText("span"),
				SimplyWallSt: e.Attr("href"),
			}
		}
	})

	c.OnScraped(func(r *colly.Response) {
		log.Print("OnScraped ", r.Request.URL)
		if page < count {
			page++
			time.Sleep(time.Duration(getRand(1, 4)) * time.Second)
			r.Request.Visit(fmt.Sprintf("https://simplywall.st/stocks/us%s?page=%d", industry, page))
		}
	})

	c.OnRequest(func(r *colly.Request) {
		log.Print("Visiting ", r.URL)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Print(err)
	})

	if err := c.Visit(fmt.Sprintf("https://simplywall.st/stocks/us%s?page=%d", industry, page)); err != nil {
		log.Print(err)
		return false
	}

	if len(result) < 10000 {
		log.Print("len(result) < 10000")
		return false
	}

	file, err := json.MarshalIndent(result, "", " ")
	if err != nil {
		log.Print(err)
		return false
	}
	err = ioutil.WriteFile("tickers.json", file, 0644)
	if err != nil {
		log.Print(err)
		return false
	}
	return true
}

func isTicker(text string) bool {
	re := regexp.MustCompile("/stocks/us/")
	loc := re.FindStringIndex(text)
	if len(loc) > 0 && loc[0] == 0 {
		if len(strings.Split(text, "/")) == 6 {
			return true
		}
	}
	return false
}

func getRand(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}
