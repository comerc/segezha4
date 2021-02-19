// package main

// import (
// 	"io/ioutil"
// 	"log"

// 	ss "github.com/comerc/segezha4/screenshot"
// )

// func main() {
// 	// linkURL := "https://marketwatch.com/investing/stock/TSLA"
// 	linkURL := "https://tipranks.com/stocks/ZM/forecast"
// 	buf := ss.MakeScreenshotForPage(linkURL, 0, 0, 0, 0)
// 	if err := ioutil.WriteFile("fullScreenshot.png", buf, 0644); err != nil {
// 		log.Fatal(err)
// 	}
// }
