package main

import (
	"io/ioutil"
	"log"

	ss "github.com/comerc/segezha4/screenshot"
)

func main() {
	// linkURL := "https://marketwatch.com/investing/stock/TSLA"
	// linkURL := "https://tipranks.com/stocks/ZM/forecast"
	linkURL := "https://www.marketbeat.com/stocks/ZM"
	buf := ss.MakeScreenshotForMarketBeat(linkURL)
	if err := ioutil.WriteFile("screenshot.png", buf, 0644); err != nil {
		log.Fatal(err)
	}
}
