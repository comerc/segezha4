package main

import (
	"context"
	"log"
	"math"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// MakePageScreenshot description
func MakePageScreenshot(linkURL string, top, height int) []byte {
	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// capture screenshot of an element
	var buf []byte

	// capture entire browser viewport, returning png with quality=90

	if err := chromedp.Run(ctx, makePageScreenshot(linkURL, top, height, 90, &buf)); err != nil {
		log.Fatal(err)
	}

	return buf
}

// makePageScreenshot takes a screenshot of the entire browser viewport.
//
// Liberally copied from puppeteer's source.
//
// Note: this will override the viewport emulation settings.
func makePageScreenshot(linkURL string, top, height int, quality int64, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(linkURL),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// get layout metrics
			_, _, contentSize, err := page.GetLayoutMetrics().Do(ctx)
			if err != nil {
				return err
			}

			width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))

			// force viewport emulation
			err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
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
				WithClip(&page.Viewport{
					X:      contentSize.X,
					Y:      float64(top), // contentSize.Y,
					Width:  contentSize.Width,
					Height: float64(height), // contentSize.Height,
					Scale:  1,
				}).Do(ctx)
			if err != nil {
				return err
			}
			return nil
		}),
	}
}

// MakeImageScreenshot description
func MakeImageScreenshot(linkURL string, width, height int) []byte {
	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// capture screenshot of an element
	var buf []byte

	// capture entire browser viewport, returning png with quality=90

	if err := chromedp.Run(ctx, makeImageScreenshot(linkURL, width, height, 90, &buf)); err != nil {
		log.Fatal(err)
	}

	return buf
}

// makeImageScreenshot takes a screenshot of the entire browser viewport.
//
// Liberally copied from puppeteer's source.
//
// Note: this will override the viewport emulation settings.
func makeImageScreenshot(linkURL string, width, height int, quality int64, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(linkURL),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// force viewport emulation
			err := emulation.SetDeviceMetricsOverride(int64(width), int64(height), 1, false).
				Do(ctx)
			if err != nil {
				return err
			}
			// capture screenshot
			*res, err = page.CaptureScreenshot().
				Do(ctx)
			if err != nil {
				return err
			}
			return nil
		}),
	}
}
