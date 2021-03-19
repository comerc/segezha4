package screenshot

import (
	"context"
	"log"
	"math"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
)

// MakeScreenshotForGuruFocus description
func MakeScreenshotForGuruFocus(linkURL string) []byte {
	// o := append(chromedp.DefaultExecAllocatorOptions[:],
	// 	// chromedp.ProxyServer("socks5://138.59.207.118:9076"),
	// 	chromedp.Flag("blink-settings", "imagesEnabled=false"),
	// )
	// ctx, cancel := chromedp.NewExecAllocator(context.Background(), o...)
	// defer cancel()
	// ctx1, cancel1 := chromedp.NewContext(ctx)
	// defer cancel1()
	ctx1, cancel1 := chromedp.NewContext(context.Background())
	defer cancel1()
	// start the browser without a timeout
	if err := chromedp.Run(ctx1); err != nil {
		log.Println(err)
		return nil
	}
	ctx2, cancel2 := context.WithTimeout(ctx1, 50*time.Second)
	defer cancel2()
	var buf []byte
	if err := chromedp.Run(ctx2, makeScreenshotForGuruFocus(linkURL, 0, 0, 0, 2042, 100, &buf)); err != nil {
		log.Println(err)
		return nil
	}
	return buf
}

// TODO: обобщить с makeScreenshotForPage
// makeScreenshotForGuruFocus takes a screenshot of the entire browser viewport.
//
// Liberally copied from puppeteer's source.
//
// Note: this will override the viewport emulation settings.
func makeScreenshotForGuruFocus(linkURL string, x, y, width, height float64, quality int64, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Emulate(device.IPadPro),
		chromedp.Navigate(linkURL),
		chromedp.WaitReady("body"),
		chromedp.SetAttributeValue("body > div.el-dialog__wrapper", "style", "display:none"),
		chromedp.SetAttributeValue("body > div.v-modal", "style", "display:none"),
		chromedp.SetAttributeValue("body > div.v-modal", "style", "display:none"),
		chromedp.Sleep(4 * time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// get layout metrics
			_, _, contentSize, err := page.GetLayoutMetrics().Do(ctx)
			if err != nil {
				return err
			}

			w, h := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))

			// force viewport emulation
			err = emulation.SetDeviceMetricsOverride(w, h, 1, false).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).
				Do(ctx)
			if err != nil {
				return err
			}

			if x == 0 {
				x = contentSize.X
			}
			if y == 0 {
				y = contentSize.Y
			}
			if width == 0 {
				width = contentSize.Width
			}
			if height == 0 {
				height = contentSize.Height
			}

			// capture screenshot
			*res, err = page.CaptureScreenshot().
				WithQuality(quality).
				WithClip(&page.Viewport{
					X:      x,
					Y:      y,
					Width:  width,
					Height: height,
					Scale:  1,
				}).Do(ctx)
			if err != nil {
				return err
			}
			return nil
		}),
	}
}
