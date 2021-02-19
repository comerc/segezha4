package screenshot

import (
	"context"
	"log"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
)

// MakeScreenshotForImage description
func MakeScreenshotForImage(linkURL string, width, height float64) []byte {
	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	var buf []byte
	if err := chromedp.Run(ctx, makeScreenshotForImage(linkURL, width, height, 100, &buf)); err != nil {
		log.Fatal(err)
	}
	return buf
}

// makeScreenshotForImage takes a screenshot of the entire browser viewport.
//
// Liberally copied from puppeteer's source.
//
// Note: this will override the viewport emulation settings.
func makeScreenshotForImage(linkURL string, width, height float64, quality int64, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Emulate(device.IPad),
		chromedp.Navigate(linkURL),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// force viewport emulation
			err := emulation.SetDeviceMetricsOverride(int64(width), int64(height), 1, false).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).
				Do(ctx)
			if err != nil {
				return err
			}
			// capture screenshot
			*res, err = page.CaptureScreenshot().
				WithQuality(quality).
				Do(ctx)
			if err != nil {
				return err
			}
			return nil
		}),
	}
}
