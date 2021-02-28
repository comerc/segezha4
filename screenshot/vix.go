package screenshot

import (
	"context"
	"log"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
)

// MakeScreenshotForVIX description
func MakeScreenshotForVIX(linkURL string) []byte {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	selChart := "body #technicalChartImage"
	var buf []byte
	if err := chromedp.Run(ctx, func() chromedp.Tasks {
		return chromedp.Tasks{
			chromedp.Emulate(device.KindleFireHDXlandscape),
			chromedp.Navigate(linkURL),
			chromedp.WaitReady("body"),
			chromedp.Sleep(1 * time.Second),
			chromedp.Screenshot(selChart, &buf, chromedp.NodeVisible),
		}
	}()); err != nil {
		log.Println(err)
	}
	return buf
}
