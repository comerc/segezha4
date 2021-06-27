package screenshot

import (
	"context"
	"log"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/comerc/segezha4/utils"
)

// MakeScreenshotForZacks description
func MakeScreenshotForZacks(linkURL string) []byte {
	o := append(chromedp.DefaultExecAllocatorOptions[:],
		// chromedp.ProxyServer("socks5://138.59.207.118:9076"),
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		chromedp.UserAgent("Mozilla/5.0"),
		chromedp.WindowSize(1024, 500),
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
	const average = 10
	ctx2, cancel2 := context.WithTimeout(ctx1, utils.GetTimeout(average))
	defer cancel2()
	selSummary := "#quote_ribbon_v2 > div.quote_rank_summary"
	selAds := "#quote_ribbon_v2 > div.quote_rank_summary > div.zer_report_box.placement_id"
	var buf []byte
	if err := chromedp.Run(ctx2, func() chromedp.Tasks {
		return chromedp.Tasks{
			network.SetBlockedURLS([]string{
				"https://*.doubleclick.net/*",
				"https://*.js*",
				"https://*.png*",
				"https://*.html*",
			}),
			// chromedp.Emulate(device.IPhoneX),
			chromedp.Navigate(linkURL),
			chromedp.WaitReady(selSummary),
			chromedp.SetAttributeValue(selAds, "style", "display:none"),
			chromedp.SetAttributeValue(selSummary, "style", "margin: 0; border: none; position: absolute; background: white; z-index: 10001;"),
			chromedp.Screenshot(selSummary, &buf, chromedp.NodeVisible),
		}
	}()); err != nil {
		log.Println(err)
		return nil
	}
	return buf
}
