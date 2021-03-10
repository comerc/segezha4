package screenshot

import (
	"bytes"
	"context"
	"image"
	"image/png"
	"log"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
)

// MakeScreenshotForFinvizIDs description
func MakeScreenshotForFinvizIDs(linkURL string) []byte {
	ctx1, cancel1 := chromedp.NewContext(context.Background())
	defer cancel1()
	// start the browser without a timeout
	if err := chromedp.Run(ctx1); err != nil {
		log.Println(err)
		return nil
	}
	ctx2, cancel2 := context.WithTimeout(ctx1, 30*time.Second)
	defer cancel2()
	var buf1, buf2 []byte
	sel1 := "body > div.content.is-index > div.container > table > tbody > tr > td table:nth-child(1)"
	sel2 := "body > div.content.is-index > div.container > table > tbody > tr > td > #homepage > table > tbody > tr > td > table"
	if err := chromedp.Run(ctx2, func() chromedp.Tasks {
		return chromedp.Tasks{
			network.SetBlockedURLS([]string{"https://dggaenaawxe8z.cloudfront.net/cmp_v2/admiral/finviz.js"}),
			chromedp.Emulate(device.KindleFireHDX),
			chromedp.Navigate(linkURL),
			chromedp.WaitReady("body #homepage"),
			// chromedp.SetAttributeValue(sel1, "style", "margin: 20px 0 0"),
			chromedp.Screenshot(sel1, &buf1, chromedp.NodeVisible),
			chromedp.SetAttributeValue(sel2, "style", "margin: 0 0 4px"),
			chromedp.Screenshot(sel2, &buf2, chromedp.NodeVisible),
		}
	}()); err != nil {
		log.Println(err)
	}
	var src image.Image
	if err := glueForCathiesArk(buf1, buf2, &src); err != nil {
		log.Println(err)
	}
	buf1 = nil
	buf2 = nil
	// resize to width 800 using Bicubic resampling
	// and preserve aspect ratio
	// res := resize.Resize(800, 0, src, resize.Bicubic)
	// encode
	out := &bytes.Buffer{}
	if err := png.Encode(out, src); err != nil {
		log.Println(err)
	}
	src = nil
	// res = nil
	return out.Bytes()
}
