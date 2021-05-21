// TODO: хостить для перехода на интерактивную версию (как и tradingview2) ???
// TODO: https://stackoverflow.com/questions/65940103/how-to-override-the-studies-of-the-tradingview-widget
// TODO: https://stackoverflow.com/questions/67433792/how-to-change-the-colors-on-the-tradingview-advanced-real-time-chart-widget-indi

package screenshot

import (
	"context"
	"log"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
	"github.com/comerc/segezha4/utils"
)

// MakeScreenshotForTradingView description
func MakeScreenshotForTradingView(linkURL string) []byte {
	ctx1, cancel1 := chromedp.NewContext(context.Background())
	defer cancel1()
	// start the browser without a timeout
	if err := chromedp.Run(ctx1); err != nil {
		log.Println(err)
		return nil
	}
	const average = 12
	ctx2, cancel2 := context.WithTimeout(ctx1, utils.GetTimeout(average))
	defer cancel2()
	var buf []byte
	container := "body > div.tradingview-widget-container"
	if err := chromedp.Run(ctx2, func() chromedp.Tasks {
		return chromedp.Tasks{
			chromedp.Emulate(device.IPadlandscape),
			chromedp.Navigate(linkURL),
			chromedp.WaitReady(container),
			chromedp.Sleep(4 * time.Second),
			chromedp.Screenshot(container, &buf, chromedp.NodeVisible),
		}
	}()); err != nil {
		log.Println(err)
		return nil
	}
	return buf
}
