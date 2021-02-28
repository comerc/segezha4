package main

import (
	"io/ioutil"
	"log"

	ss "github.com/comerc/segezha4/screenshot"
)

func _main() {

	linkURL := "https://www.barchart.com/stocks/quotes/$VIX/technical-chart/fullscreen?plot=CANDLE&volume=0&data=I:5&density=L&pricesOn=0&asPctChange=0&logscale=0&im=5&indicators=EXPMA(100);EXPMA(20);EXPMA(50);EXPMA(200);WMA(9);EXPMA(500);EXPMA(1000)&sym=$VIX&grid=1&height=625&studyheight=100"
	buf := ss.MakeScreenshotForVIX(linkURL)

	// linkURL := "https://finviz.com/quote.ashx?t=TSLA"
	// buf := ss.MakeScreenshotForFinviz(linkURL)

	// linkURL := "https://tipranks.com/stocks/ZM/forecast"
	// linkURL := "https://marketwatch.com/investing/stock/TSLA"
	// buf := ss.MakeScreenshotForMarketWatch(linkURL)

	// linkURL := "https://marketwatch.com/investing/stock/TSLA"
	// linkURL := "https://tipranks.com/stocks/ZM/forecast"
	// linkURL := "https://www.marketbeat.com/stocks/TSLA"
	// buf := ss.MakeScreenshotForMarketBeat(linkURL)
	// linkURL := "https://cathiesark.com/ark-combined-holdings-of-TSLA"
	// buf := ss.MakeScreenshotForCathiesArk(linkURL)
	// linkURL := "https://finviz.com/map.ashx?t=sec"
	// buf := ss.MakeScreenshotForFinvizMap(linkURL)
	if len(buf) == 0 {
		log.Println("exit buf == 0")
		return
	}
	if err := ioutil.WriteFile("screenshot.png", buf, 0644); err != nil {
		log.Fatal(err)
	}
}
