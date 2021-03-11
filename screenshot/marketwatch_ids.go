package screenshot

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
)

// TODO: что со шрифтами?
// TODO: sp_message_container_*

// MarketWatchTab description
type MarketWatchTab = string

// MarketWatchTab variants
const (
	MarketWatchTabUS      MarketWatchTab = "us"
	MarketWatchTabEurope  MarketWatchTab = "europe"
	MarketWatchTabAsia    MarketWatchTab = "asia"
	MarketWatchTabFX      MarketWatchTab = "fx"
	MarketWatchTabRates   MarketWatchTab = "rates"
	MarketWatchTabFutures MarketWatchTab = "futures"
	MarketWatchTabCrypto  MarketWatchTab = "crypto"
)

// MarketWatchTabs description
var MarketWatchTabs map[string]string

func init() {
	MarketWatchTabs = make(map[string]string)
	MarketWatchTabs[MarketWatchTabUS] = "https://www.marketwatch.com/markets/us"
	MarketWatchTabs[MarketWatchTabEurope] = "https://www.marketwatch.com/markets/europe-middle-east"
	MarketWatchTabs[MarketWatchTabAsia] = "https://www.marketwatch.com/markets/asia"
	MarketWatchTabs[MarketWatchTabFX] = "https://www.marketwatch.com/investing/currencies"
	MarketWatchTabs[MarketWatchTabRates] = "https://www.marketwatch.com/investing/bonds"
	MarketWatchTabs[MarketWatchTabFutures] = "https://www.marketwatch.com/investing/futures"
	MarketWatchTabs[MarketWatchTabCrypto] = "https://www.marketwatch.com/investing/cryptocurrency"
}

// MakeScreenshotForMarketWatchIDs description
func MakeScreenshotForMarketWatchIDs(linkURL string, tab MarketWatchTab) []byte {
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
			chromedp.SetAttributeValue("//body/div[starts-with(@id, 'sp_message_container_')]", "style", "display:none"),
			// chromedp.SetAttributeValue("body > #sp_message_container_450644", "style", "display:none"),
			chromedp.Click(fmt.Sprintf("//a[@href='%s']", MarketWatchTabs[tab]), chromedp.BySearch),
			chromedp.Sleep(1 * time.Second),
			chromedp.SetAttributeValue(sel, "style", "border-left:none"),
			chromedp.Screenshot(sel, &buf, chromedp.NodeVisible),
		}
	}()); err != nil {
		log.Println(err)
	}
	return buf
}
