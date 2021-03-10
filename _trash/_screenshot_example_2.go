// // Command screenshot is a chromedp example demonstrating how to take a
// // screenshot of a specific element and of the entire browser viewport.
// package main

// import (
// 	"context"
// 	"io/ioutil"
// 	"log"
// 	"math"

// 	"github.com/chromedp/cdproto/emulation"
// 	"github.com/chromedp/cdproto/page"
// 	"github.com/chromedp/chromedp"
// )

// // TODO: https://github.com/anthonynsimon/bild#crop
// // TODO: https://github.com/chromedp/examples/blob/3904eac11b/text2/main.go
// // TODO: https://habr.com/ru/post/530412/

// func main() {
// 	// create context
// 	ctx, cancel := chromedp.NewContext(context.Background())
// 	defer cancel()

// 	// capture screenshot of an element
// 	var buf []byte

// 	// if err := chromedp.Run(ctx, elementScreenshot(`https://www.gurufocus.com/stock/TAK/summary`, `#__layout`, &buf)); err != nil {
// 	// if err := chromedp.Run(ctx, elementScreenshot(`https://stockrow.com/ZM`, `#root div.capital-structure`, &buf)); err != nil {
// 	// if err := chromedp.Run(ctx, elementScreenshot(`https://finviz.com/quote.ashx?t=zm`, `body > div.content > div.container > table.snapshot-table2`, &buf)); err != nil {
// 	// if err := chromedp.Run(ctx, elementScreenshot(`https://www.marketbeat.com/stocks/NASDAQ/FB/institutional-ownership/`, `#article`, &buf)); err != nil {
// 	// 	log.Fatal(err)
// 	// }
// 	// if err := ioutil.WriteFile("elementScreenshot.png", buf, 0644); err != nil {
// 	// 	log.Fatal(err)
// 	// }

// 	// capture entire browser viewport, returning png with quality=90

// 	// if err := chromedp.Run(ctx, fullScreenshot(`https://stockrow.com/ZM`, 90, &buf)); err != nil {
// 	// if err := chromedp.Run(ctx, fullScreenshot(`https://www.gurufocus.com/stock/TAK/summary`, 90, &buf)); err != nil {
// 	// if err := chromedp.Run(ctx, fullScreenshot(`https://www.marketwatch.com/investing/stock/zm`, 90, &buf)); err != nil {
// 	// if err := chromedp.Run(ctx, fullScreenshot(`https://www.marketbeat.com/stocks/NASDAQ/FB/institutional-ownership/`, 90, &buf)); err != nil {
// 	if err := chromedp.Run(ctx, fullScreenshot(`https://finviz.com/quote.ashx?t=zm`, 90, &buf)); err != nil {
// 		log.Fatal(err)
// 	}
// 	if err := ioutil.WriteFile("fullScreenshot.png", buf, 0644); err != nil {
// 		log.Fatal(err)
// 	}

// 	log.Printf("wrote elementScreenshot.png and fullScreenshot.png")
// }

// // elementScreenshot takes a screenshot of a specific element.
// func elementScreenshot(urlstr, sel string, res *[]byte) chromedp.Tasks {
// 	return chromedp.Tasks{
// 		// emulation.SetUserAgentOverride("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3830.0 Safari/537.36"),
// 		chromedp.Navigate(urlstr),
// 		// chromedp.Sleep(8 * time.Second),
// 		// chromedp.WaitReady(`body > div > footer`),
// 		// chromedp.WaitVisible(`#optinform-modal`),
// 		// chromedp.Click(`#optinform-modal a`, chromedp.NodeVisible),
// 		// chromedp.WaitReady(sel, chromedp.ByID),
// 		chromedp.WaitVisible(sel),
// 		chromedp.Screenshot(sel, res, chromedp.NodeVisible),
// 	}
// }

// // fullScreenshot takes a screenshot of the entire browser viewport.
// //
// // Liberally copied from puppeteer's source.
// //
// // Note: this will override the viewport emulation settings.
// func fullScreenshot(urlstr string, quality int64, res *[]byte) chromedp.Tasks {
// 	return chromedp.Tasks{
// 		chromedp.Navigate(urlstr),
// 		// chromedp.Sleep(8 * time.Second),
// 		// chromedp.Click(`#root div.close-modal`, chromedp.NodeVisible),

// 		// chromedp.WaitVisible(`#optinform-modal`),
// 		// chromedp.Click(`#optinform-modal a`, chromedp.NodeVisible),
// 		chromedp.ActionFunc(func(ctx context.Context) error {
// 			// get layout metrics
// 			_, _, contentSize, err := page.GetLayoutMetrics().Do(ctx)
// 			if err != nil {
// 				return err
// 			}

// 			width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))

// 			// force viewport emulation
// 			err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
// 				WithScreenOrientation(&emulation.ScreenOrientation{
// 					Type:  emulation.OrientationTypePortraitPrimary,
// 					Angle: 0,
// 				}).
// 				Do(ctx)
// 			if err != nil {
// 				return err
// 			}

// 			// capture screenshot
// 			*res, err = page.CaptureScreenshot().
// 				WithQuality(quality).
// 				WithClip(&page.Viewport{
// 					X:      contentSize.X,
// 					Y:      contentSize.Y, // + 710,
// 					Width:  contentSize.Width,
// 					Height: contentSize.Height, // 540,
// 					Scale:  1,
// 				}).Do(ctx)
// 			if err != nil {
// 				return err
// 			}
// 			return nil
// 		}),
// 	}
// }
