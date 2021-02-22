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
	var buf11, buf12 []byte
	if err := takeScreenshot(ctx, "#liInsiderTrades > a", "#insiderChart", &buf11, &buf12); err != nil {
		log.Println(err)
	}
	var buf21, buf22 []byte
	if err := takeScreenshot(ctx, "#liInstutionalOwnership > a", "#SECChart", &buf21, &buf22); err != nil {
		log.Println(err)
	}
	var src1 image.Image
	if err := glue(buf12, buf11, &src1); err != nil {
		log.Println(err)
	}
	buf11 = nil
	buf12 = nil
	var src2 image.Image
	if err := glue(buf22, buf21, &src2); err != nil {
		log.Println(err)
	}
	buf21 = nil
	buf22 = nil
	var src image.Image
	if src1 == nil && src2 == nil {
		return nil
	}
	if src1 != nil && src2 != nil {
		glueImages(src1, src2, &src)
		src1 = nil
		src2 = nil
	} else {
		if src1 != nil {
			src = src1
		}
		if src2 != nil {
			src = src2
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
	src1 = nil
	src2 = nil
	src = nil
	res = nil
	// var opt jpeg.Options
	// opt.Quality = 85
	// jpeg.Encode(out, rgba, &opt)
	return out.Bytes()
}

func takeScreenshot(ctx context.Context, linkSel, chartSel interface{}, titleRes, chartRes *[]byte) error {
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
	if err := chromedp.Run(ctx, func() chromedp.Tasks {
		return chromedp.Tasks{
			chromedp.Click(linkSel, chromedp.NodeVisible),
			chromedp.WaitReady(`body > div > footer`),
			chromedp.WaitVisible(chartSel),
		}
	}()); err != nil {
		return err
	}
	titleSel := "#article > #form1 > #cphPrimaryContent_pnlCompany > #shareableArticle > div:nth-child(2) > div > div"
	if err := chromedp.Run(ctx, func() chromedp.Tasks {
		return chromedp.Tasks{
			chromedp.SetAttributeValue(titleSel, "style", "padding:8px"),
			chromedp.Screenshot(titleSel, titleRes, chromedp.NodeVisible),
		}
	}()); err != nil {
		return err
	}
	sel := fmt.Sprintf("%v > #svg > #yTextGroup > g.footnote", chartSel)
	if err := chromedp.Run(ctx, func() chromedp.Tasks {
		return chromedp.Tasks{
			chromedp.Nodes(sel, &nodes, chromedp.AtLeast(0)),
		}
	}()); err != nil {
		return err
	}
	if len(nodes) == 0 {
		return nil
	}
	if err := chromedp.Run(ctx, func() chromedp.Tasks {
		return chromedp.Tasks{
			chromedp.SetAttributeValue(sel, "style", "display:none"),
			chromedp.Screenshot(chartSel, chartRes, chromedp.NodeVisible),
		}
	}()); err != nil {
		return err
	}
	return nil
}

func glue(buf1, buf2 []byte, src *image.Image) error {
	if len(buf1) == 0 && len(buf2) == 0 {
		return nil
	}
	if len(buf1) > 0 && len(buf2) > 0 {
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
	} else {
		if len(buf1) > 0 {
			img1, _, err := image.Decode(bytes.NewReader(buf1))
			if err != nil {
				return err
			}
			buf1 = nil
			*src = img1
		}
		if len(buf2) > 0 {
			img2, _, err := image.Decode(bytes.NewReader(buf2))
			if err != nil {
				return err
			}
			buf2 = nil
			*src = img2
		}
	}
	return nil
}

func glueImages(img1, img2 image.Image, src *image.Image) error {
	//starting position of the second image (bottom left)
	sp2 := image.Point{0, img1.Bounds().Dy()}
	//new rectangle for the second image
	r2 := image.Rectangle{sp2, sp2.Add(img2.Bounds().Size())}
	//rectangle for the big image
	r1 := image.Rectangle{image.Point{0, 0}, r2.Max}
	rgba := image.NewRGBA(r1)
	draw.Draw(rgba, img1.Bounds(), img1, image.Point{0, 0}, draw.Src)
	draw.Draw(rgba, r2, img2, image.Point{0, 0}, draw.Src)
	*src = rgba
	return nil
}
