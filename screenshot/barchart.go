package screenshot

import (
	"context"
	"log"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
)

// MakeScreenshotForBarChart description
func MakeScreenshotForBarChart(linkURL string) []byte {
	ctx0, cancel0 := chromedp.NewRemoteAllocator(context.Background(), getWebSocketDebuggerUrl())
	defer cancel0()
	ctx1, cancel1 := chromedp.NewContext(ctx0)
	defer cancel1()
	// start the browser without a timeout
	if err := chromedp.Run(ctx1); err != nil {
		log.Println(err)
		return nil
	}
	ctx2, cancel2 := context.WithTimeout(ctx1, 50*time.Second)
	defer cancel2()
	selChart := "body #technicalChartImage"
	var buf []byte
	if err := chromedp.Run(ctx2, func() chromedp.Tasks {
		return chromedp.Tasks{
			chromedp.Emulate(device.KindleFireHDX),
			chromedp.Navigate(linkURL),
			chromedp.WaitReady("body"),
			// chromedp.Sleep(1 * time.Second),
			chromedp.Sleep(4 * time.Second),
			chromedp.Click("//*[text()='Accept all']", chromedp.BySearch),
			chromedp.Screenshot(selChart, &buf, chromedp.NodeVisible),
		}
	}()); err != nil {
		log.Println(err)
		return nil
	}
	return buf
}
