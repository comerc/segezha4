package screenshot

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"
	"log"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
)

// MakeScreenshotForTipRanks description
func MakeScreenshotForTipRanks(symbol string) []byte {
	// o := append(chromedp.DefaultExecAllocatorOptions[:],
	// 	chromedp.ProxyServer("socks5://138.59.207.118:9076"),
	// 	// chromedp.Flag("blink-settings", "imagesEnabled=false"),
	// )
	// ctx, cancel := chromedp.NewExecAllocator(context.Background(), o...)
	// defer cancel()
	// ctx1, cancel1 := chromedp.NewContext(ctx)
	// defer cancel1()
	ctx1, cancel1 := chromedp.NewContext(context.Background())
	defer cancel1()
	// start the browser without a timeout
	if err := chromedp.Run(ctx1); err != nil {
		log.Println(err)
		return nil
	}
	ctx2, cancel2 := context.WithTimeout(ctx1, 50*time.Second)
	defer cancel2()
	sel1 := "body div.client-components-stock-research-smart-score-style__rank"
	sel2 := "body div.client-components-stock-research-analysts-style__analystTopPart"
	sel3 := "body div.client-components-stock-research-individual-investors-style__topSection"
	sel4 := "body div.client-components-stock-research-bloggers-bloggerStyles__generalOpinions"
	var buf1, buf2, buf3, buf4 []byte
	if err := chromedp.Run(ctx2, func() chromedp.Tasks {
		return chromedp.Tasks{
			chromedp.Emulate(device.IPadlandscape),
			chromedp.Navigate(fmt.Sprintf("https://www.tipranks.com/stocks/%s/forecast", symbol)),
			chromedp.WaitReady("body > #app "),
			chromedp.Sleep(4 * time.Second),
			chromedp.SetAttributeValue("body > #app > div > div > div.tr-app", "style", "display:none"),
			chromedp.Click("body nav > a:nth-child(1)", chromedp.BySearch),
			chromedp.SetAttributeValue(sel1, "style", "margin: 10px 0"),
			chromedp.ActionFunc(screenshotWithoutPopups(sel1, &buf1)),
			chromedp.Click("body nav > a:nth-child(2)", chromedp.BySearch),
			chromedp.SetAttributeValue(sel2, "style", "margin: 0 0 10px"),
			chromedp.ActionFunc(screenshotWithoutPopups(sel2, &buf2)),
			chromedp.Click("body nav > a:nth-child(3)", chromedp.BySearch),
			chromedp.SetAttributeValue(sel3, "style", "margin: 0 0 10px"),
			chromedp.ActionFunc(screenshotWithoutPopups(sel3, &buf3)),
			chromedp.Click("body nav > a:nth-child(4)", chromedp.BySearch),
			chromedp.SetAttributeValue(sel4, "style", "margin: 0 0 10px"),
			chromedp.ActionFunc(screenshotWithoutPopups(sel4, &buf4)),
		}
	}()); err != nil {
		log.Println(err)
		return nil
	}
	var src image.Image
	if err := glueForTipRanks(buf1, buf2, buf3, buf4, &src); err != nil {
		log.Println(err)
		return nil
	}
	buf1, buf2, buf3, buf4 = nil, nil, nil, nil
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

func screenshotWithoutPopups(sel string, buf *[]byte) func(context.Context) error {
	var fn func(context.Context) error
	fn = func(ctx context.Context) error {
		if err := chromedp.Screenshot(sel, buf, chromedp.NodeVisible).Do(ctx); err != nil {
			return err
		}
		var isPopup1, isPopup2 bool
		if err := hidePopup(ctx, "body #gtm_popup_blocker_iframe", &isPopup1); err != nil {
			return err
		}
		if err := hidePopup(ctx, "body > #popup-ios-modal-v4", &isPopup2); err != nil {
			return err
		}
		if isPopup1 || isPopup2 {
			return fn(ctx)
		}
		return nil
	}
	return fn
}

func hidePopup(ctx context.Context, sel string, isPopup *bool) error {
	var nodes []*cdp.Node
	if err := chromedp.Nodes(sel, &nodes, chromedp.AtLeast(0)).Do(ctx); err != nil {
		return err
	}
	if len(nodes) == 0 {
		return nil
	}
	var ok = false
	var value string
	if err := chromedp.AttributeValue(sel, "style", &value, &ok).Do(ctx); err != nil {
		return err
	}
	if ok && value == "display:none" {
		return nil
	}
	if err := chromedp.SetAttributeValue(sel, "style", "display:none").Do(ctx); err != nil {
		return err
	}
	*isPopup = true
	return nil
}

func glueForTipRanks(buf1, buf2, buf3, buf4 []byte, src *image.Image) error {
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
	img4, _, err := image.Decode(bytes.NewReader(buf4))
	if err != nil {
		return err
	}
	buf4 = nil
	glueImages(img1, img2, src)
	glueImages(*src, img3, src)
	glueImages(*src, img4, src)
	img1 = nil
	img2 = nil
	img3 = nil
	img4 = nil
	return nil
}
