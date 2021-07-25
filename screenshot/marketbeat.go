package screenshot

import (
	"bytes"
	"context"
	"image"
	"image/png"
	"log"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
	"github.com/comerc/segezha4/utils"
)

// MakeScreenshotForMarketBeat description
func MakeScreenshotForMarketBeat(linkURL string) []byte {
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
	const average = 25 * 2 // *2 !!
	ctx2, cancel2 := context.WithTimeout(ctx1, utils.GetTimeout(average))
	defer cancel2()
	if err := chromedp.Run(ctx2, func() chromedp.Tasks {
		return chromedp.Tasks{
			network.SetBlockedURLS([]string{
				// "*.ashx*",
				"https://cdn.onesignal.com/sdks/OneSignalSDK.js",
				"https://www.marketbeat.com/scripts/modal/mb-modernpopup.js?v=*",
				"https://pagead2.googlesyndication.com/pagead/js/adsbygoogle.js",
				"https://www.statcounter.com/counter/counter.js",
				"https://www.googletagmanager.com/gtag/js?id=AW-1050184556",
				"https://instant.page/*",
				"https://www.marketbeat.com/cdn-cgi/scripts/*/cloudflare-static/email-decode.min.js",
				// "https://cdnjs.cloudflare.com/ajax/libs/d3/*/d3.min.js",
				// "https://www.marketbeat.com/Scripts/charts/dv-charts.js?v=*",
				"https://cdnjs.cloudflare.com/ajax/libs/hellojs/*/hello.all.js",
				"https://ajax.googleapis.com/ajax/libs/jquery/*/jquery.min.js",
				"https://accounts.google.com/gsi/client",
				"https://cdnjs.cloudflare.com/ajax/libs/popper.js/*/umd/popper.min.js",
				"https://stackpath.bootstrapcdn.com/bootstrap/*/js/bootstrap.min.js",
				"https://code.jquery.com/ui/*/jquery-ui.min.js",
				"https://www.marketbeat.com/scripts/masterscripts9.js?v=*",
				"https://www.statcounter.com/counter/counter.js",
				"https://www.google-analytics.com/analytics.js",
				// "chrome-extension://*/js/inject.js",
				"*.aspx*", "*.json*", "*.png*"}),
			chromedp.Emulate(device.KindleFireHDX),
			chromedp.Navigate(linkURL),
			chromedp.WaitReady("body footer"),
			// v2
			// chromedp.WaitVisible("#optinform-container"),
			// chromedp.SetAttributeValue("#optinform-container", "style", "display:none"),
			// v1
			// chromedp.WaitVisible("#optinform-modal a"),
			// chromedp.Click("#optinform-modal a", chromedp.NodeVisible),
		}
	}()); err != nil {
		log.Println(err)
		return nil
	}
	var buf11, buf12 []byte
	if err := takeScreenshotForMarketBeat(ctx2, "#liInsiderTrades > a", "#insiderChart", &buf11, &buf12); err != nil {
		log.Println(err)
		return nil
	}
	var buf21, buf22 []byte
	if err := takeScreenshotForMarketBeat(ctx2, "#liInstutionalOwnership > a", "#SECChart", &buf21, &buf22); err != nil {
		log.Println(err)
		return nil
	}
	var src1 image.Image
	if err := glueForMarketBeat(buf12, buf11, &src1); err != nil {
		log.Println(err)
		return nil
	}
	buf11 = nil
	buf12 = nil
	var src2 image.Image
	if err := glueForMarketBeat(buf22, buf21, &src2); err != nil {
		log.Println(err)
		return nil
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
	// res := resize.Resize(800, 0, src, resize.Bicubic)
	// encode
	out := &bytes.Buffer{}
	if err := png.Encode(out, src); err != nil {
		log.Println(err)
		return nil
	}
	src1 = nil
	src2 = nil
	src = nil
	// res = nil
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
	pnlAdSenseHeader := "#pnlAdSenseHeader"
	if err := chromedp.Run(ctx, func() chromedp.Tasks {
		return chromedp.Tasks{
			chromedp.Click(linkSel, chromedp.NodeVisible),
			chromedp.WaitReady("body footer"),
			chromedp.Sleep(1 * time.Second),
			// chromedp.WaitVisible(chartSel),
			// chromedp.ScrollIntoView(chartSel),
			chromedp.SetAttributeValue(pnlAdSenseHeader, "style", "display:none"),
			chromedp.WaitNotVisible(pnlAdSenseHeader),
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
			chromedp.SetAttributeValue(titleSel, "style", "padding:8px"),
			chromedp.Screenshot(titleSel, titleRes, chromedp.NodeVisible),
		}
	}()); err != nil {
		return err
	}
	// sel := fmt.Sprintf("%v > #svg > #yTextGroup > g.footnote", chartSel)
	// if err := chromedp.Run(ctx, func() chromedp.Tasks {
	// 	return chromedp.Tasks{
	// 		chromedp.Nodes(sel, &nodes, chromedp.AtLeast(0)),
	// 	}
	// }()); err != nil {
	// 	return err
	// }
	// if len(nodes) == 0 {
	// 	return nil
	// }
	if err := chromedp.Run(ctx, func() chromedp.Tasks {
		return chromedp.Tasks{
			// chromedp.SetAttributeValue(sel, "style", "display:none"),
			// chromedp.WaitNotVisible(sel),
			// chromedp.Screenshot(chartSel, chartRes, chromedp.NodeVisible),
			chromedp.Screenshot("#wrapper", chartRes, chromedp.NodeVisible),
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
