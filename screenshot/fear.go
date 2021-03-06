package screenshot

import (
	"context"
	"log"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
	"github.com/comerc/segezha4/utils"
)

// MakeScreenshotForFear description
func MakeScreenshotForFear(linkURL string) []byte {
	ctx1, cancel1 := chromedp.NewContext(context.Background())
	defer cancel1()
	// start the browser without a timeout
	if err := chromedp.Run(ctx1); err != nil {
		log.Println(err)
		return nil
	}
	const average = 14
	ctx2, cancel2 := context.WithTimeout(ctx1, utils.GetTimeout(average))
	defer cancel2()
	var buf []byte
	body := "body > #cnnBody"
	sel := body + " div.mod-quoteinfo.feargreed"
	if err := chromedp.Run(ctx2, func() chromedp.Tasks {
		return chromedp.Tasks{
			network.SetBlockedURLS([]string{
				// "*.js*",
				// "*.json*",
				// "*.png",
				"*.svg",
				"*.gif",
				"*.jpg",
				"*.mp4",
				"*.html",
				"*doubleclick.net*",
				"*google.com*",
				"*googlesyndication.com*",
				"*sharethis.com*",
			}),
			chromedp.Emulate(device.IPadPro),
			chromedp.Navigate(linkURL),
			chromedp.WaitReady(body),
			chromedp.Sleep(4 * time.Second),
			chromedp.SetAttributeValue("body > #onetrust-consent-sdk", "style", "display:none"),
			chromedp.SetAttributeValue("body > #cnnHeader", "style", "display:none"),
			chromedp.SetAttributeValue("body > #adBanner", "style", "display:none"),
			chromedp.SetAttributeValue(body, "style", "margin-top:0"),
			chromedp.SetAttributeValue(body+" div.indicatorHeading", "style", "display:none"),
			chromedp.SetAttributeValue(body+" div.indicatorContainer", "style", "display:none"),
			chromedp.SetAttributeValue(sel, "style", "border:none"),
			chromedp.ScrollIntoView(body),
			chromedp.Sleep(2 * time.Second),
			chromedp.Screenshot(sel, &buf, chromedp.NodeVisible),
		}
	}()); err != nil {
		log.Println(err)
		return nil
	}
	return buf
}
