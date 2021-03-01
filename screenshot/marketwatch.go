package screenshot

import (
	"context"
	"log"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
)

// MakeScreenshotForMarketWatch description
func MakeScreenshotForMarketWatch(linkURL string) []byte {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	var buf []byte
	if err := chromedp.Run(ctx, func() chromedp.Tasks {
		return chromedp.Tasks{
			chromedp.Emulate(device.IPad),
			chromedp.Navigate(linkURL),
			chromedp.WaitReady(`body > footer`),
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
	return buf
}
