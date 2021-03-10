package screenshot

import (
	"bytes"
	"context"
	"image"
	"image/png"
	"log"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
)

// MakeScreenshotForCathiesArk description
func MakeScreenshotForCathiesArk(linkURL string) []byte {
	ctx1, cancel1 := chromedp.NewContext(context.Background())
	defer cancel1()
	// start the browser without a timeout
	if err := chromedp.Run(ctx1); err != nil {
		log.Println(err)
		return nil
	}
	ctx2, cancel2 := context.WithTimeout(ctx1, 40*time.Second)
	defer cancel2()
	var buf1, buf2 []byte
	sel0 := "body main > div:nth-child(2) > div:nth-child(2)"
	sel1 := "body header"
	sel2 := "body main > div:nth-child(1)"
	sel3 := "body main div.ant-row.sectionContainer___plkQX:nth-child(3) > div > div > div.recharts-responsive-container > div.recharts-wrapper"
	if err := chromedp.Run(ctx2, func() chromedp.Tasks {
		return chromedp.Tasks{
			chromedp.Emulate(device.IPadlandscape),
			chromedp.Navigate(linkURL),
			chromedp.WaitReady("body"),
			chromedp.SetAttributeValue("body main", "style", ""),
			chromedp.SetAttributeValue(sel0, "style", ""),
			chromedp.SetAttributeValue(sel0+" > div", "style", "padding: 40px 40px 0; flex: 0 0 100%; max-width: 100%;"),
			chromedp.Screenshot(sel0, &buf1, chromedp.NodeVisible),
			chromedp.SetAttributeValue(sel1, "style", "display:none"),
			chromedp.SetAttributeValue(sel2, "style", "display:none"),
			// chromedp.SetAttributeValue(sel3+" > svg > g:nth-child(4) > g", "style", "display:none"),
			// chromedp.SetAttributeValue(sel3+" > div.recharts-legend-wrapper", "style", "display:none"),
			chromedp.Sleep(4 * time.Second),
			// TODO: убирать, если успел появиться
			// chromedp.SetAttributeValue("body > div > div.ant-notification", "style", "display:none"),
			chromedp.Screenshot(sel3, &buf2, chromedp.NodeVisible),
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

func glueForCathiesArk(buf1, buf2 []byte, src *image.Image) error {
	img1, _, err := image.Decode(bytes.NewReader(buf1))
	if err != nil {
		return err
	}
	buf1 = nil
	img2, _, err := image.Decode(bytes.NewReader(buf2))
	if err != nil {
		return err
	}
	buf2 = nil
	glueImages(img1, img2, src)
	return nil
}
