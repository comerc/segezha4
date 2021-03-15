package screenshot

import (
	"context"
	"log"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
)

// MakeScreenshotForMarketWatch description
func MakeScreenshotForMarketWatch(linkURL string) []byte {
	o := append(chromedp.DefaultExecAllocatorOptions[:],
		// chromedp.ProxyServer("socks5://138.59.207.118:9076"),
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
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
	ctx2, cancel2 := context.WithTimeout(ctx1, 40*time.Second)
	defer cancel2()
	// var s string
	var buf []byte
	if err := chromedp.Run(ctx2, func() chromedp.Tasks {
		return chromedp.Tasks{
			chromedp.Emulate(device.IPad),
			chromedp.Navigate(linkURL),
			chromedp.WaitReady("body > footer"),
			chromedp.Sleep(4 * time.Second),
			chromedp.SetAttributeValue("//body/div[starts-with(@id, 'sp_message_container_')]", "style", "display:none"),
			// chromedp.SetAttributeValue("body > #sp_message_container_450644", "style", "display:none"),
			chromedp.SetAttributeValue("body > div.container.container--body > div.region.region--intraday > div.column.column--full.quote__nav", "style", "display:none"),
			chromedp.SetAttributeValue("body > div.container.container--body > div.region.region--intraday > div.column.column--full > div.element.element--company > div.row", "style", "display:none"),
			chromedp.SetAttributeValue("body > div.container.container--body > div.region.region--intraday > div.column.column--full > div.element.element--company > div.row > div.quote-actions", "style", "display:none"),
			chromedp.SetAttributeValue("body > div.container.container--body > div.region.region--intraday > div.column.column--primary > mw-chart.element.element--chart > label.toggle--chart", "style", "display:none"),
			chromedp.SetAttributeValue("body > div.container.container--body > div.region.region--intraday > div.column.column--primary > mw-chart.element.element--chart > div.chart__options", "style", "display:none"),
			chromedp.SetAttributeValue("body > div.container.container--trending", "style", "display:none"),
			chromedp.Screenshot("body > div.container.container--body > div.region.region--intraday", &buf, chromedp.NodeVisible),
		}
	}()); err != nil {
		log.Println(err)
	}
	// d1 := []byte(s)
	// if err := ioutil.WriteFile("/tmp/dat_mw.html", d1, 0644); err != nil {
	// 	log.Println(err)
	// }
	return buf
}
