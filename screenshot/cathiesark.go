package screenshot

import (
	"context"
	"log"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
)

// MakeScreenshotForCathiesArk description
func MakeScreenshotForCathiesArk(linkURL string) []byte {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	sel1 := "body header"
	sel2 := "body main > div:nth-child(1)"
	sel3 := "body main div.ant-row.sectionContainer___plkQX > div > div"
	var buf []byte
	if err := chromedp.Run(ctx, func() chromedp.Tasks {
		return chromedp.Tasks{
			chromedp.Emulate(device.IPadProlandscape),
			chromedp.Navigate(linkURL),
			chromedp.WaitReady("body"),
			chromedp.Sleep(4 * time.Second),
			chromedp.SetAttributeValue(sel1, "style", "display:none"),
			chromedp.SetAttributeValue(sel2, "style", "display:none"),
			chromedp.SetAttributeValue(sel3+" svg > g:nth-child(4) > g", "style", "display:none"),
			chromedp.SetAttributeValue(sel3+" div.recharts-legend-wrapper", "style", "display:none"),
			chromedp.Screenshot(sel3, &buf, chromedp.NodeVisible),
		}
	}()); err != nil {
		log.Println(err)
	}
	return buf
}
