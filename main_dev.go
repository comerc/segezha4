package main

import (
	"io/ioutil"
	"log"

	ss "github.com/comerc/segezha4/screenshot"
	"github.com/comerc/segezha4/utils"
	"github.com/joho/godotenv"
)

func _main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}
	utils.InitTimeoutFactor()

	// linkURL := "https://money.cnn.com/data/fear-and-greed/"
	// buf := ss.MakeScreenshotForFear(linkURL)

	// linkURL := "https://www.barchart.com/stocks/quotes/$VIX/technical-chart/fullscreen?plot=CANDLE&volume=0&data=I:5&density=L&pricesOn=0&asPctChange=0&logscale=0&im=5&indicators=EXPMA(100);EXPMA(20);EXPMA(50);EXPMA(200);WMA(9);EXPMA(500);EXPMA(1000)&sym=$VIX&grid=1&height=625&studyheight=100"
	// buf := ss.MakeScreenshotForVIX(linkURL)

	// linkURL := "https://finviz.com/quote.ashx?t=ZM"
	// buf := ss.MakeScreenshotForFinviz(linkURL)

	// linkURL := "https://marketwatch.com/"
	// buf := ss.MakeScreenshotForMarketWatchIDs(linkURL, ss.MarketWatchHrefCrypto)

	// linkURL := "https://www.marketwatch.com/"
	// buf := ss.MakeScreenshotForMarketWatchIDs(linkURL, ss.MarketWatchTabUS)

	// linkURL := "https://marketwatch.com/investing/stock/TSLA"
	// buf := ss.MakeScreenshotForMarketWatch(linkURL)

	// linkURL := "https://www.gurufocus.com/stock/amd/summary#"
	// buf := ss.MakeScreenshotForGuruFocus(linkURL)

	// linkURL := "https://www.gurufocus.com/stock/irbt/summary#"
	// buf := ss.MakeScreenshotForPage(linkURL, 0, 0, 0, 2042)

	// linkURL := "https://tipranks.com/stocks/life/forecast"
	// buf := ss.MakeScreenshotForTipRanks(linkURL)

	// linkURL := "https://tipranks.com/stocks/life/stock-analysis"
	// buf := ss.MakeScreenshotForTipRanks2(linkURL)

	// linkURL := "https://marketwatch.com/investing/stock/TSLA"
	// linkURL := "https://www.marketbeat.com/stocks/TSLA"
	// buf := ss.MakeScreenshotForMarketBeat(linkURL)

	// linkURL := "https://cathiesark.com/ark-combined-holdings-of-tsla"
	// buf := ss.MakeScreenshotForCathiesArk(linkURL)

	// linkURL := "https://finviz.com/map.ashx?t=sec"
	// buf := ss.MakeScreenshotForFinvizMap(linkURL)

	// linkURL := "https://finviz.com/"
	// buf := ss.MakeScreenshotForFinvizBB(linkURL)

	// path, _ := os.Getwd()
	// path = filepath.Join(path, "assets/tradingview.html")
	// symbol := "MU"
	// interval := "4H"
	// linkURL := fmt.Sprintf("file://%s?%s:%s", path, symbol, interval)
	// buf := ss.MakeScreenshotForTradingView(linkURL)

	// path, _ := os.Getwd()
	// path = filepath.Join(path, "assets/tradingview2.html")
	// symbol := "MU"
	// interval1 := "4H"
	// interval2 := "1H"
	// linkURL := fmt.Sprintf("file://%s?%s:%s:%s", path, symbol, interval1, interval2)
	// buf := ss.MakeScreenshotForTradingView2(linkURL)

	// path, _ := os.Getwd()
	// path = filepath.Join(path, "assets/bestday.html")
	// now := time.Now()
	// day := fmt.Sprintf("%02d-%02d", now.Month(), now.Day())
	// linkURL := fmt.Sprintf("file://%s?%s", path, day)
	// buf := ss.MakeScreenshotForBestDay(linkURL)

	// linkURL := "https://zacks.com/stock/quote/tsla"
	// buf := ss.MakeScreenshotForZacks(linkURL)

	linkURL := "https://simplywall.st/stocks/us/automobiles/nasdaq-tsla/tesla"
	_, buf := ss.MakeScreenshotForSimplyWallSt(linkURL)

	if len(buf) == 0 {
		log.Println("exit buf == 0")
		return
	}
	if err := ioutil.WriteFile("_screenshot.png", buf, 0644); err != nil {
		log.Fatal(err)
	}
}
