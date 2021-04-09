package main

// import (
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"strings"
// 	"sync"
// 	"time"

// 	ss "github.com/comerc/segezha4/screenshot"
// 	"github.com/comerc/segezha4/utils"
// )

// func _main() {

// 	// linkURL := "https://money.cnn.com/data/fear-and-greed/"
// 	// buf := ss.MakeScreenshotForFear(linkURL)

// 	// linkURL := "https://www.barchart.com/stocks/quotes/$VIX/technical-chart/fullscreen?plot=CANDLE&volume=0&data=I:5&density=L&pricesOn=0&asPctChange=0&logscale=0&im=5&indicators=EXPMA(100);EXPMA(20);EXPMA(50);EXPMA(200);WMA(9);EXPMA(500);EXPMA(1000)&sym=$VIX&grid=1&height=625&studyheight=100"
// 	// buf := ss.MakeScreenshotForVIX(linkURL)

// 	// linkURL := "https://finviz.com/quote.ashx?t=ZM"
// 	// buf := ss.MakeScreenshotForFinviz(linkURL)

// 	// linkURL := "https://finviz.com/"
// 	// buf := ss.MakeScreenshotForFinvizIDs(linkURL)

// 	// linkURL := "https://marketwatch.com/"
// 	// buf := ss.MakeScreenshotForMarketWatchIDs(linkURL, ss.MarketWatchHrefCrypto)

// 	// linkURL := "https://www.marketwatch.com/"
// 	// buf := ss.MakeScreenshotForMarketWatchIDs(linkURL, ss.MarketWatchTabUS)

// 	// linkURL := "https://marketwatch.com/investing/stock/TSLA"
// 	// buf := ss.MakeScreenshotForMarketWatch(linkURL)

// 	// linkURL := "https://www.gurufocus.com/stock/amd/summary#"
// 	// buf := ss.MakeScreenshotForGuruFocus(linkURL)

// 	// linkURL := "https://www.gurufocus.com/stock/irbt/summary#"
// 	// buf := ss.MakeScreenshotForPage(linkURL, 0, 0, 0, 2042)

// 	// linkURL := "https://tipranks.com/stocks/ZM/forecast"
// 	// buf := ss.MakeScreenshotForTipRanks(linkURL)

// 	// linkURL := "https://marketwatch.com/investing/stock/TSLA"
// 	// linkURL := "https://www.marketbeat.com/stocks/TSLA"
// 	// buf := ss.MakeScreenshotForMarketBeat(linkURL)

// 	// linkURL := "https://cathiesark.com/ark-combined-holdings-of-flir"
// 	// buf := ss.MakeScreenshotForCathiesArk(linkURL)

// 	// linkURL := "https://finviz.com/map.ashx?t=sec"
// 	// buf := ss.MakeScreenshotForFinvizMap(linkURL)
// 	// if len(buf) == 0 {
// 	// 	log.Println("exit buf == 0")
// 	// 	return
// 	// }
// 	// if err := ioutil.WriteFile("screenshot.png", buf, 0644); err != nil {
// 	// 	log.Fatal(err)
// 	// }

// 	// o := append(chromedp.DefaultExecAllocatorOptions[:]) // chromedp.ProxyServer("socks5://138.59.207.118:9076"),
// 	// // chromedp.Flag("blink-settings", "imagesEnabled=false"),
// 	// // chromedp.DisableGPU,

// 	// ctx, cancel := chromedp.NewExecAllocator(context.Background(), o...)
// 	// defer cancel()

// 	// /info finviz.com TSLA ZM BYND
// 	// articleCase := GetExactArticleCase("tipranks.com")
// 	// if articleCase == nil {
// 	// 	sendText(b, m.Chat.ID, "Invalid command", false)
// 	// 	return
// 	// }
// 	dirtySymbols := []string{"tsla", "fb", "zm", "bynd"}

// 	symbols := normalizeSymbols(dirtySymbols)

// 	callbacks := make([]getWhat, len(symbols))
// 	for i, symbol := range symbols {
// 		callbacks[i] = closeWhat(symbol)
// 	}

// 	_ = callbacks
// 	sendBatch(false, callbacks)
// }

// type ParallelResult struct {
// 	buf        []byte
// 	isReceived bool
// 	isSent     bool
// }

// type getWhat func() []byte

// func closeWhat(symbol string) getWhat {
// 	return func() []byte {
// 		linkURL := fmt.Sprintf("https://tipranks.com/stocks/%s/forecast", symbol)
// 		return ss.MakeScreenshotForTipRanks(linkURL)
// 	}
// }

// func sendBatch(isPrivateChat bool, callbacks []getWhat) {
// 	defer elapsed("start")()
// 	done := make(chan bool)
// 	results := make([]ParallelResult, len(callbacks))
// 	threads := 3
// 	fmt.Println("threads", threads)
// 	if threads == 0 {
// 		threads = 1
// 	}
// 	var tokens = make(chan struct{}, threads) // ограничение количества горутин
// 	var mu sync.Mutex
// 	receivedCount := 0
// 	for i, cb := range callbacks {
// 		tokens <- struct{}{} // захват маркера
// 		go func(i int, cb getWhat) {
// 			defer elapsed(fmt.Sprintf("screenshot%d.png", i))()
// 			buf := cb()
// 			<-tokens // освобождение маркера
// 			{
// 				mu.Lock()
// 				defer mu.Unlock()
// 				results[i] = ParallelResult{
// 					buf:        buf,
// 					isReceived: true,
// 				}
// 				receivedCount = receivedCount + 1
// 				if receivedCount == len(callbacks) {
// 					sendAllReceived2(isPrivateChat, results, len(results))
// 					done <- true
// 				} else {
// 					isAllPreviosReceived := true
// 					for _, r := range results[:i] {
// 						if !r.isReceived {
// 							isAllPreviosReceived = false
// 							break
// 						}
// 					}
// 					if isAllPreviosReceived {
// 						sendAllReceived2(isPrivateChat, results, i+1)
// 					}
// 				}
// 			}
// 		}(i, cb)
// 	}
// 	<-done
// }

// var lastSend2 = time.Now().AddDate(0, 0, -1)

// func sendAllReceived2(isPrivateChat bool, results []ParallelResult, l int) {
// 	// fmt.Println("sendAllReceived")
// 	for i, r := range results[:l] {
// 		func(i int, r ParallelResult) {
// 			if !r.isSent {
// 				if !isPrivateChat {
// 					// your bot will not be able to send more than 20 messages per minute to the same group.
// 					diff := time.Since(lastSend2)
// 					if diff < 4*time.Second {
// 						time.Sleep(4 * time.Second)
// 						lastSend2 = time.Now()
// 					}
// 				}
// 				if len(r.buf) == 0 {
// 					log.Println("exit buf == 0 " + fmt.Sprintf("screenshot%d.png", i))
// 					return
// 				}
// 				if err := ioutil.WriteFile(fmt.Sprintf("screenshot%d.png", i), r.buf, 0644); err != nil {
// 					log.Fatal(err)
// 				}
// 				results[i].isSent = true
// 			}
// 		}(i, r)
// 	}
// }

// func elapsed(what string) func() {
// 	start := time.Now()
// 	return func() {
// 		fmt.Printf("%s took %v\n", what, time.Since(start))
// 	}
// }

// func isNotFoundTicker(symbol string) bool {
// 	// TODO: реализация пополнения тикеров
// 	// ticker := GetExactTicker(symbol)
// 	// return ticker == nil
// 	return false
// }

// func normalizeSymbols(symbols []string) []string {
// 	result := make([]string, 0)
// 	for _, symbol := range symbols {
// 		if strings.HasPrefix(symbol, "#") || strings.HasPrefix(symbol, "$") {
// 			symbol = symbol[1:]
// 		}
// 		if utils.Contains(result, strings.ToUpper(symbol)) {
// 			continue
// 		}
// 		result = append(result, strings.ToUpper(symbol))
// 	}
// 	return result
// }

func workaround() {

}
