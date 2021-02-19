package screenshot

import (
	"context"
	"log"
	"math"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
)

// MakeScreenshotForPage description
func MakeScreenshotForPage(linkURL string, x, y, width, height float64) []byte {
	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// capture screenshot of an element
	var buf []byte

	// capture entire browser viewport, returning png with quality=90

	if err := chromedp.Run(ctx, makeScreenshotForPage(linkURL, x, y, width, height, 100, &buf)); err != nil {
		log.Fatal(err)
	}

	return buf
}

// makeScreenshotForPage takes a screenshot of the entire browser viewport.
//
// Liberally copied from puppeteer's source.
//
// Note: this will override the viewport emulation settings.
func makeScreenshotForPage(linkURL string, x, y, width, height float64, quality int64, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Emulate(device.KindleFireHDX),
		chromedp.Navigate(linkURL),
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
