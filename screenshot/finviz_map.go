package screenshot

import (
	"context"
	"log"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
)

// MakeScreenshotForFinvizMap description
func MakeScreenshotForFinvizMap(linkURL string) []byte {
	ctx0, cancel0 := chromedp.NewRemoteAllocator(context.Background(), getWebSocketDebuggerUrl())
	defer cancel0()
	ctx1, cancel1 := chromedp.NewContext(ctx0)
	defer cancel1()
	// start the browser without a timeout
	if err := chromedp.Run(ctx1); err != nil {
		log.Println(err)
		return nil
	}
	ctx2, cancel2 := context.WithTimeout(ctx1, 50*time.Second)
	defer cancel2()
	// TODO: отключить картинки
	selHeader := "body > table.header"
	selNavbar := "body > table.navbar"
	selView := "body > div.content.map > div.container > div.view"
	selChart := "body > div.content.map > div.container > div.content-view-map > #map > #body > div > div > canvas.chart"
	selFooter := "body > div.content.map > div.container > div.content-view-map > #map > #body > div > div:nth-child(2)"
	var buf []byte
	if err := chromedp.Run(ctx2, func() chromedp.Tasks {
		return chromedp.Tasks{
			chromedp.Emulate(device.KindleFireHDXlandscape),
			chromedp.Navigate(linkURL),
			chromedp.WaitReady("body"),
			chromedp.SetAttributeValue(selHeader, "style", "display:none"),
			chromedp.SetAttributeValue(selNavbar, "style", "display:none"),
			chromedp.SetAttributeValue(selView, "style", "display:none"),
			chromedp.SetAttributeValue(selChart, "style", "margin:6px 0 0 2px"),
			chromedp.SetAttributeValue(selFooter, "style", "display:none"),
			chromedp.Screenshot(selChart, &buf, chromedp.NodeVisible),
		}
	}()); err != nil {
		log.Println(err)
		return nil
	}
	return buf
}
