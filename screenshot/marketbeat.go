package screenshot

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"
	"log"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
	"github.com/nfnt/resize"
)

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
	if err := takeScreenshotForMarketBeat(ctx, "#liInsiderTrades > a", "#insiderChart", &buf11, &buf12); err != nil {
		log.Println(err)
	}
	var buf21, buf22 []byte
	if err := takeScreenshotForMarketBeat(ctx, "#liInstutionalOwnership > a", "#SECChart", &buf21, &buf22); err != nil {
		log.Println(err)
	}
	var src1 image.Image
	if err := glueForMarketBeat(buf12, buf11, &src1); err != nil {
		log.Println(err)
	}
	buf11 = nil
	buf12 = nil
	var src2 image.Image
	if err := glueForMarketBeat(buf22, buf21, &src2); err != nil {
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

func takeScreenshotForMarketBeat(ctx context.Context, linkSel, chartSel interface{}, titleRes, chartRes *[]byte) error {
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
			chromedp.WaitReady(chartSel),
		}
	}()); err != nil {
		return err
	}

	selBar := "body > #mb-bar"
	if err := chromedp.Run(ctx, func() chromedp.Tasks {
		return chromedp.Tasks{
			chromedp.Nodes(selBar, &nodes, chromedp.AtLeast(0)),
		}
	}()); err != nil {
		return err
	}
	if len(nodes) == 1 {
		if err := chromedp.Run(ctx, func() chromedp.Tasks {
			return chromedp.Tasks{
				chromedp.SetAttributeValue(selBar, "style", "display:none"),
				chromedp.WaitNotVisible(selBar),
			}
		}()); err != nil {
			return err
		}
	}

	titleSel := "#article > #form1 > #cphPrimaryContent_pnlCompany > #shareableArticle > div:nth-child(2) > div > div"
	if err := chromedp.Run(ctx, func() chromedp.Tasks {
		return chromedp.Tasks{
			chromedp.WaitReady(titleSel),
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
			chromedp.WaitNotVisible(sel),
			chromedp.Screenshot(chartSel, chartRes, chromedp.NodeVisible),
		}
	}()); err != nil {
		return err
	}
	return nil
}

func glueForMarketBeat(buf1, buf2 []byte, src *image.Image) error {
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
