package screenshot

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
	"github.com/nfnt/resize"
)

func init() {
	// image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
}

// MakeScreenshotForMarketBeat description
func MakeScreenshotForMarketBeat(linkURL string) []byte {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	if err := chromedp.Run(ctx, func() chromedp.Tasks {
		return chromedp.Tasks{
			chromedp.Emulate(device.KindleFireHDX),
			chromedp.Navigate(linkURL),
			chromedp.WaitReady(`body > div > footer`),
			chromedp.WaitVisible("#optinform-modal a"),
			chromedp.Click("#optinform-modal a", chromedp.NodeVisible),
		}
	}()); err != nil {
		log.Println(err)
	}
	var buf1 []byte
	if err := makeScreenshot(ctx, "#liInsiderTrades > a", "#insiderChart", &buf1); err != nil {
		log.Println(err)
	}
	var buf2 []byte
	if err := makeScreenshot(ctx, "#liInstutionalOwnership > a", "#SECChart", &buf2); err != nil {
		log.Println(err)
	}
	if len(buf1) == 0 && len(buf2) == 0 {
		return nil
	}
	var src image.Image
	if len(buf1) > 0 && len(buf2) > 0 {
		img1, _, err := image.Decode(bytes.NewReader(buf1))
		if err != nil {
			log.Println(err)
		}
		img2, _, err := image.Decode(bytes.NewReader(buf2))
		if err != nil {
			log.Println(err)
		}
		//starting position of the second image (bottom left)
		sp2 := image.Point{0, img1.Bounds().Dy()}
		//new rectangle for the second image
		r2 := image.Rectangle{sp2, sp2.Add(img2.Bounds().Size())}
		//rectangle for the big image
		r1 := image.Rectangle{image.Point{0, 0}, r2.Max}
		rgba := image.NewRGBA(r1)
		draw.Draw(rgba, img1.Bounds(), img1, image.Point{0, 0}, draw.Src)
		draw.Draw(rgba, r2, img2, image.Point{0, 0}, draw.Src)
		src = rgba
	} else {
		if len(buf1) > 0 {
			img1, _, err := image.Decode(bytes.NewReader(buf1))
			if err != nil {
				log.Println(err)
			}
			src = img1
		}
		if len(buf2) > 0 {
			img2, _, err := image.Decode(bytes.NewReader(buf2))
			if err != nil {
				log.Println(err)
			}
			src = img2
		}
	}
	// resize to width 800 using Bicubic resampling
	// and preserve aspect ratio
	res := resize.Resize(800, 0, src, resize.Bicubic)
	// encode
	out := &bytes.Buffer{}
	if err := png.Encode(out, res); err != nil {
		log.Println(err)
	}
	// var opt jpeg.Options
	// opt.Quality = 85
	// jpeg.Encode(out, rgba, &opt)
	return out.Bytes()
}

func makeScreenshot(ctx context.Context, linkSel, chartSel interface{}, res *[]byte) error {
	var nodes []*cdp.Node
	if err := chromedp.Run(ctx, func() chromedp.Tasks {
		return chromedp.Tasks{
			chromedp.Nodes(linkSel, &nodes, chromedp.AtLeast(0)),
		}
	}()); err != nil {
		return err
	}
	if len(nodes) == 0 {
		return nil
	}
	sel := fmt.Sprintf("%v > #svg > #yTextGroup > g.footnote", chartSel)
	if err := chromedp.Run(ctx, func() chromedp.Tasks {
		return chromedp.Tasks{
			chromedp.Click(linkSel, chromedp.NodeVisible),
			chromedp.WaitReady(`body > div > footer`),
			chromedp.WaitVisible(chartSel),
			chromedp.SetAttributeValue(sel, "style", "display:none"),
			chromedp.Screenshot(chartSel, res, chromedp.NodeVisible),
		}
	}()); err != nil {
		return err
	}
	return nil
}
