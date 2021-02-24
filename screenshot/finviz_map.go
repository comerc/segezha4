package screenshot

import (
	"context"
	"log"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
)

// MakeScreenshotForFinvizMap description
func MakeScreenshotForFinvizMap(linkURL string) []byte {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	selHeader := "body > table.header"
	selNavbar := "body > table.navbar"
	selView := "body > div.content.map > div.container > div.view"
	selChart := "body > div.content.map > div.container > div.content-view-map > #map > #body > div > div > canvas.chart"
	var buf []byte
	if err := chromedp.Run(ctx, func() chromedp.Tasks {
		return chromedp.Tasks{
			chromedp.Emulate(device.KindleFireHDXlandscape),
			chromedp.Navigate(linkURL),
			chromedp.WaitReady("body"),
			chromedp.SetAttributeValue(selHeader, "style", "display:none"),
			chromedp.SetAttributeValue(selNavbar, "style", "display:none"),
			chromedp.SetAttributeValue(selView, "style", "display:none"),
			chromedp.Screenshot(selChart, &buf, chromedp.NodeVisible),
		}
	}()); err != nil {
		log.Println(err)
	}
	return buf
}
