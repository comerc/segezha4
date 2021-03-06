package screenshot

import (
	"bytes"
	"context"
	"image"
	"image/png"
	"log"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
	"github.com/comerc/segezha4/utils"
)

// MakeScreenshotForFinviz description
func MakeScreenshotForFinviz(linkURL string) []byte {
	o := append(chromedp.DefaultExecAllocatorOptions[:],
		// chromedp.ProxyServer("socks5://138.59.207.118:9076"),
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
	)
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), o...)
	defer cancel()
	ctx1, cancel1 := chromedp.NewContext(ctx)
	defer cancel1()
	// ctx1, cancel1 := chromedp.NewContext(context.Background())
	// defer cancel1()
	// start the browser without a timeout
	if err := chromedp.Run(ctx1); err != nil {
		log.Println(err)
		return nil
	}
	const average = 8
	ctx2, cancel2 := context.WithTimeout(ctx1, utils.GetTimeout(average))
	defer cancel2()
	if err := chromedp.Run(ctx2, func() chromedp.Tasks {
		return chromedp.Tasks{
			network.SetBlockedURLS([]string{"https://dggaenaawxe8z.cloudfront.net/cmp_v2/admiral/finviz.js"}),
			chromedp.Emulate(device.KindleFireHDX),
			chromedp.Navigate(linkURL),
			chromedp.WaitReady("body"),
			// chromedp.Sleep(4 * time.Second),
			// chromedp.Click("//*[text()='Accept all']", chromedp.BySearch),
		}
	}()); err != nil {
		log.Println(err)
		return nil
	}
	var buf1, buf2, buf3 []byte
	if err := takeScreenshotForFinviz(ctx2, &buf1, &buf2, &buf3); err != nil {
		log.Println(err)
	}
	if len(buf1) == 0 {
		return nil
	}
	var src image.Image
	if err := glueForFinviz(buf1, buf2, buf3, &src); err != nil {
		log.Println(err)
		return nil
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
		return nil
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
	img3, _, err := image.Decode(bytes.NewReader(buf3))
	if err != nil {
		return err
	}
	buf3 = nil
	glueImages(img1, img2, src)
	glueImages(*src, img3, src)
	img1 = nil
	img2 = nil
	img3 = nil
	return nil
}

func takeScreenshotForFinviz(ctx context.Context, buf1, buf2, buf3 *[]byte) error {
	selChart := "body > div > #app > #chart > #charts"
	selTitleContainer := "body > div.content > div.fv-container > table:nth-child(1) > tbody > tr > td"
	selTitle := "body > div.content > div.fv-container > table:nth-child(1) table:nth-child(1)"
	selTable := "body > div.content > div.fv-container > table:nth-child(2)"
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
			chromedp.SetAttributeValue(selTitleContainer, "style", "padding:0"),
			chromedp.Screenshot(selTitle, buf2, chromedp.NodeVisible),
			chromedp.SetAttributeValue(selTable, "style", "margin-top:16px"),
			chromedp.Screenshot(selTable, buf3, chromedp.NodeVisible),
		}
	}()); err != nil {
		return err
	}
	return nil
}
