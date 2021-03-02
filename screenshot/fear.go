package screenshot

import (
	"context"
	"log"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
)

// MakeScreenshotForFear description
func MakeScreenshotForFear(linkURL string) []byte {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	var buf []byte
	body := "body > #cnnBody"
	sel := body + " div.mod-quoteinfo.feargreed"
	if err := chromedp.Run(ctx, func() chromedp.Tasks {
		return chromedp.Tasks{
			chromedp.Emulate(device.IPadPro),
			chromedp.Navigate(linkURL),
			chromedp.WaitReady(body),
			chromedp.SetAttributeValue(body+" div.indicatorHeading", "style", "display:none"),
			chromedp.SetAttributeValue(body+" div.indicatorContainer", "style", "display:none"),
			chromedp.SetAttributeValue(sel, "style", "border:none"),
			chromedp.Screenshot(sel, &buf, chromedp.NodeVisible),
		}
	}()); err != nil {
		log.Println(err)
	}
	return buf
}
