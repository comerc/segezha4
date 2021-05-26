package screenshot

import (
	"context"
	"log"

	"github.com/chromedp/chromedp"
	"github.com/comerc/segezha4/utils"
)

// MakeScreenshotForFinvizMap description
func MakeScreenshotForFinvizMap(linkURL string) []byte {
	o := append(chromedp.DefaultExecAllocatorOptions[:],
		// chromedp.ProxyServer("socks5://138.59.207.118:9076"),
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		chromedp.UserAgent("Mozilla/5.0"),
		chromedp.WindowSize(1600, 1200),
	)
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), o...)
	defer cancel()
	ctx1, cancel1 := chromedp.NewContext(ctx)
	defer cancel1()
	// ctx1, cancel1 := chromedp.NewContext(context.Background())
	// defer cancel1()
	// start the browser without a timeout
	if err := chromedp.Run(ctx1); err != nil {
		log.Println(err)
		return nil
	}
	const average = 9
	ctx2, cancel2 := context.WithTimeout(ctx1, utils.GetTimeout(average))
	defer cancel2()
	selHeader := "body > table.header"
	selNavbar := "body > table.navbar"
	selView := "body > div.content.map > div.fv-container > div.view"
	selChart := "body > div.content.map > div.fv-container > div.content-view-map > #map > #body > div > div > canvas.chart"
	selFooter := "body > div.content.map > div.fv-container > div.content-view-map > #map > #body > div > div:nth-child(2)"
	var buf []byte
	if err := chromedp.Run(ctx2, func() chromedp.Tasks {
		return chromedp.Tasks{
			// chromedp.Emulate(device.KindleFireHDXlandscape),
			chromedp.Navigate(linkURL),
			chromedp.WaitReady("body"),
			chromedp.SetAttributeValue(selHeader, "style", "display:none"),
			chromedp.SetAttributeValue(selNavbar, "style", "display:none"),
			chromedp.SetAttributeValue(selView, "style", "display:none"),
			// chromedp.SetAttributeValue(selChart, "style", "margin:2px 0 0 2px; width: 788px;height: 438px;"),
			chromedp.SetAttributeValue(selChart, "style", "margin:2px 1px 1px 2px; width: 1211px; height: 672px;"),
			chromedp.SetAttributeValue(selFooter, "style", "display:none"),
			chromedp.Screenshot(selChart, &buf, chromedp.NodeVisible),
		}
	}()); err != nil {
		log.Println(err)
		return nil
	}
	return buf
}
