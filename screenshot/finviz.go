package screenshot

import (
	"bytes"
	"context"
	"image"
	"image/png"
	"log"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
)

// MakeScreenshotForFinviz description
func MakeScreenshotForFinviz(linkURL string) []byte {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	if err := chromedp.Run(ctx, func() chromedp.Tasks {
		return chromedp.Tasks{
			chromedp.Emulate(device.KindleFireHDX),
			chromedp.Navigate(linkURL),
			chromedp.WaitReady("body"),
		}
	}()); err != nil {
		log.Println(err)
	}
	var buf1, buf2, buf3 []byte
	if err := takeScreenshotForFinviz(ctx, &buf1, &buf2, &buf3); err != nil {
		log.Println(err)
	}
	if len(buf1) == 0 {
		return nil
	}
	var src image.Image
	if err := glueForFinviz(buf1, buf2, buf3, &src); err != nil {
		log.Println(err)
	}
	buf1 = nil
	buf2 = nil
	buf3 = nil
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

func glueForFinviz(buf1, buf2, buf3 []byte, src *image.Image) error {
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
	var img3 image.Image
	img4, _, err := image.Decode(bytes.NewReader(buf3))
	if err != nil {
		return err
	}
	buf3 = nil
	// TODO: glueImages3
	glueImages(img1, img2, &img3)
	glueImages(img3, img4, src)
	img1 = nil
	img2 = nil
	img3 = nil
	img4 = nil
	return nil
}

func takeScreenshotForFinviz(ctx context.Context, buf1, buf2, buf3 *[]byte) error {
	selChart := "body > div > #app > #chart > #charts"
	selTitle := "body > div.content > div.container table.fullview-title > tbody"
	selTable := "body > div.content > div.container > table.snapshot-table2"
	var nodes []*cdp.Node
	if err := chromedp.Run(ctx, func() chromedp.Tasks {
		return chromedp.Tasks{
			chromedp.Nodes(selChart, &nodes, chromedp.AtLeast(0)),
		}
	}()); err != nil {
		return err
	}
	if len(nodes) == 0 {
		return nil
	}
	if err := chromedp.Run(ctx, func() chromedp.Tasks {
		return chromedp.Tasks{
			chromedp.SetAttributeValue(selChart, "style", "margin-bottom:10px"),
			chromedp.Screenshot(selChart, buf1, chromedp.NodeVisible),
			chromedp.Screenshot(selTitle, buf2, chromedp.NodeVisible),
			chromedp.SetAttributeValue(selTable, "style", "margin-top:16px"),
			chromedp.Screenshot(selTable, buf3, chromedp.NodeVisible),
		}
	}()); err != nil {
		return err
	}
	return nil
}
