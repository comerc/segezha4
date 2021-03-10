package screenshot

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
)

// MarketWatchTab description
type MarketWatchTab = string

// MarketWatchTab variants
const (
	MarketWatchTabUS      MarketWatchTab = "https://www.marketwatch.com/markets/us"
	MarketWatchTabEurope  MarketWatchTab = "https://www.marketwatch.com/markets/europe-middle-east"
	MarketWatchTabAsia    MarketWatchTab = "https://www.marketwatch.com/markets/asia"
	MarketWatchTabFX      MarketWatchTab = "https://www.marketwatch.com/investing/currencies"
	MarketWatchTabRates   MarketWatchTab = "https://www.marketwatch.com/investing/bonds"
	MarketWatchTabFutures MarketWatchTab = "https://www.marketwatch.com/investing/futures"
	MarketWatchTabCrypto  MarketWatchTab = "https://www.marketwatch.com/investing/cryptocurrency"
)

// MakeScreenshotForMarketWatchIDs description
func MakeScreenshotForMarketWatchIDs(linkURL string, tab MarketWatchTab) []byte {
	ctx1, cancel1 := chromedp.NewContext(context.Background())
	defer cancel1()
	// start the browser without a timeout
	if err := chromedp.Run(ctx1); err != nil {
		log.Println(err)
		return nil
	}
	ctx2, cancel2 := context.WithTimeout(ctx1, 30*time.Second)
	defer cancel2()
	// sel := "body > section > div.region.region--full.masthead-elements > div.column.column--full.max-markets > div.element.element--markets.desktop > div.markets--desktop"
	sel := "body div.markets__table > table"
	var buf []byte
	if err := chromedp.Run(ctx2, func() chromedp.Tasks {
		return chromedp.Tasks{
			chromedp.Emulate(device.IPad),
			chromedp.Navigate(linkURL),
			chromedp.WaitReady("body > footer"),
			chromedp.Sleep(4 * time.Second),
			chromedp.SetAttributeValue("body > #sp_message_container_413120", "style", "display:none"),
			chromedp.Click(fmt.Sprintf("//a[@href='%s']", tab), chromedp.BySearch),
			chromedp.Sleep(1 * time.Second),
			chromedp.SetAttributeValue(sel, "style", "border-left:none"),
			chromedp.Screenshot(sel, &buf, chromedp.NodeVisible),
		}
	}()); err != nil {
		log.Println(err)
	}
	return buf
}
