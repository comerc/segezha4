package main

import (
	"io/ioutil"
	"log"

	ss "github.com/comerc/segezha4/screenshot"
)

func _main() {

	// linkURL := "https://finviz.com/quote.ashx?t=TSLA"
	// buf := ss.MakeScreenshotForFinviz(linkURL)

	// linkURL := "https://marketwatch.com/investing/stock/BABA"
	// // linkURL := "https://tipranks.com/stocks/ZM/forecast"
	// buf := ss.MakeScreenshotForMarketWatch(linkURL)

	// linkURL := "https://marketwatch.com/investing/stock/TSLA"
	// linkURL := "https://tipranks.com/stocks/ZM/forecast"
	// linkURL := "https://www.marketbeat.com/stocks/TSLA"
	// buf := ss.MakeScreenshotForMarketBeat(linkURL)
	linkURL := "https://cathiesark.com/ark-combined-holdings-of-TSLA"
	buf := ss.MakeScreenshotForCathiesArk(linkURL)
	if len(buf) == 0 {
		log.Println("exit buf == 0")
		return
	}
	if err := ioutil.WriteFile("screenshot.png", buf, 0644); err != nil {
		log.Fatal(err)
	}
}
