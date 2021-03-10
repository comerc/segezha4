package screenshot

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
)

// MarketWatchHref description
type MarketWatchHref = string

// MarketWatchHref variants
const (
	MarketWatchHrefUS      MarketWatchHref = "https://www.marketwatch.com/markets/us"
	MarketWatchHrefEurope  MarketWatchHref = "https://www.marketwatch.com/markets/europe-middle-east"
	MarketWatchHrefAsia    MarketWatchHref = "https://www.marketwatch.com/markets/asia"
	MarketWatchHrefFX      MarketWatchHref = "https://www.marketwatch.com/investing/currencies"
	MarketWatchHrefRates   MarketWatchHref = "https://www.marketwatch.com/investing/bonds"
	MarketWatchHrefFutures MarketWatchHref = "https://www.marketwatch.com/investing/futures"
	MarketWatchHrefCrypto  MarketWatchHref = "https://www.marketwatch.com/investing/cryptocurrency"
)

// MarketWatchTab struct
type MarketWatchTab struct {
	name string
	href MarketWatchHref
}

// MarketWatchTabs slice
var MarketWatchTabs = []MarketWatchTab{
	{name: "us", href: MarketWatchHrefUS},
	{name: "europe", href: MarketWatchHrefEurope},
	{name: "asia", href: MarketWatchHrefAsia},
	{name: "fx", href: MarketWatchHrefFX},
	{name: "rates", href: MarketWatchHrefRates},
	{name: "futures", href: MarketWatchHrefFutures},
	{name: "crypto", href: MarketWatchHrefCrypto},
}

// elements := make(map[string]string)

// MakeScreenshotForMarketWatchIDs description
func MakeScreenshotForMarketWatchIDs(linkURL string, tabHref MarketWatchHref) []byte {
	ctx1, cancel1 := chromedp.NewContext(context.Background())
	defer cancel1()
	// start the browser without a timeout
	if err := chromedp.Run(ctx1); err != nil {
		log.Println(err)
		return nil
	}
	ctx2, cancel2 := context.WithTimeout(ctx1, 40*time.Second)
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
			// chromedp.SetAttributeValue("body > #sp_message_container_413120", "style", "display:none"),
			chromedp.Click(fmt.Sprintf("//a[@href='%s']", tabHref), chromedp.BySearch),
			chromedp.Sleep(1 * time.Second),
			chromedp.SetAttributeValue(sel, "style", "border-left:none"),
			chromedp.Screenshot(sel, &buf, chromedp.NodeVisible),
		}
	}()); err != nil {
		log.Println(err)
	}
	return buf
}
