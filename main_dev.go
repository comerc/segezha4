package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"sync"
	"time"

	ss "github.com/comerc/segezha4/screenshot"
)

type Result struct {
	buf        []byte
	isReceived bool
	isSent     bool
}

func main() {

	// linkURL := "https://money.cnn.com/data/fear-and-greed/"
	// buf := ss.MakeScreenshotForFear(linkURL)

	// linkURL := "https://www.barchart.com/stocks/quotes/$VIX/technical-chart/fullscreen?plot=CANDLE&volume=0&data=I:5&density=L&pricesOn=0&asPctChange=0&logscale=0&im=5&indicators=EXPMA(100);EXPMA(20);EXPMA(50);EXPMA(200);WMA(9);EXPMA(500);EXPMA(1000)&sym=$VIX&grid=1&height=625&studyheight=100"
	// buf := ss.MakeScreenshotForVIX(linkURL)

	// linkURL := "https://finviz.com/quote.ashx?t=ZM"
	// buf := ss.MakeScreenshotForFinviz(linkURL)

	// linkURL := "https://finviz.com/"
	// buf := ss.MakeScreenshotForFinvizIDs(linkURL)

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

	// linkURL := "https://tipranks.com/stocks/ZM/forecast"
	// buf := ss.MakeScreenshotForTipRanks("tsla")

	// linkURL := "https://marketwatch.com/investing/stock/TSLA"
	// linkURL := "https://www.marketbeat.com/stocks/TSLA"
	// buf := ss.MakeScreenshotForMarketBeat(linkURL)

	// linkURL := "https://cathiesark.com/ark-combined-holdings-of-flir"
	// buf := ss.MakeScreenshotForCathiesArk(linkURL)

	// linkURL := "https://finviz.com/map.ashx?t=sec"
	// buf := ss.MakeScreenshotForFinvizMap(linkURL)
	// if len(buf) == 0 {
	// 	log.Println("exit buf == 0")
	// 	return
	// }
	// if err := ioutil.WriteFile("screenshot.png", buf, 0644); err != nil {
	// 	log.Fatal(err)
	// }

	a := []string{"tsla", "fb", "zm", "bynd", "ge", "gm", "mu", "aa"}

	cbs := make([]Callback, len(a))

	for i, symbol := range a {
		func(i int, symbol string) {
			cb := func() []byte {
				return ss.MakeScreenshotForTipRanks(symbol)
			}
			cbs[i] = cb
		}(i, symbol)
	}

	sendBatch(cbs)
}

type Callback func() []byte

func sendBatch(cbs []Callback) {
	done := make(chan bool)
	defer elapsed("start")()
	results := make([]Result, len(cbs))

	// opts := append(
	// 	chromedp.DefaultExecAllocatorOptions[:],
	// 	// select all the elements after the third element
	// 	// chromedp.DefaultExecAllocatorOptions[3:],
	// 	// chromedp.NoFirstRun,
	// 	// chromedp.NoDefaultBrowserCheck,
	// 	chromedp.DisableGPU,
	// )

	// ctx0, cancel0 := chromedp.NewExecAllocator(context.Background(), opts...)
	// // ctx0, cancel0 := chromedp.NewExecAllocator(context.Background())
	// defer cancel0()

	// ctx1, cancel1 := chromedp.NewContext(ctx0)
	// defer cancel1()
	// // start the browser without a timeout
	// if err := chromedp.Run(ctx1); err != nil {
	// 	log.Fatalln(err)
	// }

	var tokens = make(chan struct{}, 4) // ограничение количества горутин
	var mu sync.Mutex
	receivedCount := 0
	for i, cb := range cbs {
		tokens <- struct{}{} // захват маркера
		go func(i int, cb Callback) {
			// defer elapsed(fmt.Sprintf("screenshot%d.png", i))()
			// buf := ss.MakeScreenshotForTipRanks(t)
			buf := cb()
			<-tokens // освобождение маркера
			{
				mu.Lock()
				defer mu.Unlock()
				results[i] = Result{
					buf:        buf,
					isReceived: true,
				}
				receivedCount = receivedCount + 1
				if receivedCount == len(cbs) {
					sendAllReceived(results, len(results))
					done <- true
				} else {
					isAllPreviosReceived := true
					for _, r := range results[:i] {
						if !r.isReceived {
							isAllPreviosReceived = false
							break
						}
					}
					if isAllPreviosReceived {
						sendAllReceived(results, i+1)
					}
				}
			}
		}(i, cb)
	}
	<-done
}

func sendAllReceived(results []Result, l int) {
	for i, r := range results[:l] {
		if !r.isSent {
			results[i].isSent = true
			if len(r.buf) == 0 {
				log.Println("buf == 0")
			} else if err := ioutil.WriteFile(fmt.Sprintf("screenshot%d.png", i), r.buf, 0644); err != nil {
				log.Println(err)
			}
		}
	}
}

func elapsed(what string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", what, time.Since(start))
	}
}
