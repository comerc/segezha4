package screenshot

import (
	"context"
	"log"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
)

// MakeScreenshotForFinviz description
func MakeScreenshotForFinviz(linkURL string) []byte {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	var buf []byte
	if err := chromedp.Run(ctx, func() chromedp.Tasks {
		return chromedp.Tasks{
			chromedp.Emulate(device.IPad),
			chromedp.Navigate(linkURL),
			chromedp.WaitReady("body"),
			chromedp.Screenshot("#app > #chart > #charts", &buf, chromedp.NodeVisible),
		}
	}()); err != nil {
		log.Println(err)
	}
	return buf
}
