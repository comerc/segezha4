// Стратегия «Три экрана Элдера» https://smart-lab.ru/blog/568328.php https://alpari.com/ru/beginner/articles/tri-ekrana-eldera/

package screenshot

import (
	"context"
	"log"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
	"github.com/comerc/segezha4/utils"
)

// MakeScreenshotForTradingView2 description
func MakeScreenshotForTradingView2(linkURL string) []byte {
	ctx1, cancel1 := chromedp.NewContext(context.Background())
	defer cancel1()
	// start the browser without a timeout
	if err := chromedp.Run(ctx1); err != nil {
		log.Println(err)
		return nil
	}
	const average = 18
	ctx2, cancel2 := context.WithTimeout(ctx1, utils.GetTimeout(average))
	defer cancel2()
	var buf []byte
	container := "body > div.tradingview-widget-container"
	if err := chromedp.Run(ctx2, func() chromedp.Tasks {
		return chromedp.Tasks{
			chromedp.Emulate(device.IPadPro),
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
