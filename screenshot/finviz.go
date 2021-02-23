package screenshot

import (
	"bytes"
	"context"
	"image"
	"image/png"
	"log"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
	"github.com/nfnt/resize"
)

// MakeScreenshotForFinviz description
func MakeScreenshotForFinviz(linkURL string) []byte {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	selChart := "body > div > #app > #chart > #charts"
	selTable := "body > div.content > div.container > table.snapshot-table2"
	var buf1, buf2 []byte
	if err := chromedp.Run(ctx, func() chromedp.Tasks {
		return chromedp.Tasks{
			chromedp.Emulate(device.IPadPro),
			chromedp.Navigate(linkURL),
			chromedp.WaitReady("body"),
			chromedp.SetAttributeValue(selChart, "style", "padding:20px"),
			chromedp.Screenshot(selChart, &buf1, chromedp.NodeVisible),
			chromedp.SetAttributeValue(selTable, "style", "padding:20px"),
			chromedp.Screenshot(selTable, &buf2, chromedp.NodeVisible),
		}
	}()); err != nil {
		log.Println(err)
	}

	var src image.Image
	if err := glueForFinviz(buf1, buf2, &src); err != nil {
		log.Println(err)
	}
	buf1 = nil
	buf2 = nil
	// resize to width 800 using Bicubic resampling
	// and preserve aspect ratio
	res := resize.Resize(800, 0, src, resize.Bicubic)
	// encode
	out := &bytes.Buffer{}
	if err := png.Encode(out, res); err != nil {
		log.Println(err)
	}
	src = nil
	res = nil
	return out.Bytes()
}

func glueForFinviz(buf1, buf2 []byte, src *image.Image) error {
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
