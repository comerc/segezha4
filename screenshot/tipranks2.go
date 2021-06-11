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

// MakeScreenshotForTipRanks2 description
func MakeScreenshotForTipRanks2(linkURL string) []byte {
	// o := append(chromedp.DefaultExecAllocatorOptions[:],
	// 	// chromedp.ProxyServer("socks5://138.59.207.118:9076"),
	// 	// chromedp.Flag("blink-settings", "imagesEnabled=false"),
	// 	// chromedp.DisableGPU,
	// )
	// ctx, cancel := chromedp.NewExecAllocator(context.Background(), o...)
	// defer cancel()
	// ctx1, cancel1 := chromedp.NewContext(ctx)
	// defer cancel1()
	// ctx1, cancel1 := chromedp.NewContext(context.Background())
	// defer cancel1()
	// start the browser without a timeout
	// if err := chromedp.Run(ctx1); err != nil {
	// 	log.Println(err)
	// 	return nil
	// }
	ctx0 := context.Background()
	ctx1, cancel1 := chromedp.NewContext(ctx0)
	defer cancel1()
	// start the browser without a timeout
	if err := chromedp.Run(ctx1, func() chromedp.Tasks {
		return chromedp.Tasks{
			network.SetBlockedURLS([]string{
				"https://blog.tipranks.com/*",
				"https://randomuser.me/*",
				"/new-images/stock-research/banner/*",
			}),
		}
	}()); err != nil {
		log.Println(err)
		return nil
	}
	const average = 11
	ctx2, cancel2 := context.WithTimeout(ctx1, utils.GetTimeout(average))
	defer cancel2()
	selHeader := "#root > div:nth-child(2) > div.shadowheader.ipad_shadownone"
	sel1 := "#tr-stock-page-content > div.maxW1200.grow1.flexc__.flexc__.displayflex > div.minW80.z1.flexr__f.maxW1200.mobile_maxWparent > div.tr-box-ui.flexc__.w12.displayflex.minHauto.z0.mb7.mobile_px0.mobile_pr0.mobile_pl0.mobile_w12 > div.flexc__.mt3.bgwhite.displayflex.border1.borderColorwhite-8.shadow1.positionrelative.grow1 > div.w12.displayflex.positionrelative.grow1.ipadpro_pl0.ipadpro_pr0.desktop_flexc__ > div > div.bgwhite.flexcb_.mt0.displayflex.desktop_pl4 > div.flexccc.w12.displayflex > div"
	sel2 := "#tr-stock-page-content > div.maxW1200.grow1.flexc__.flexc__.displayflex > div.minW80.z1.flexr__f.maxW1200.mobile_maxWparent > div.tr-box-ui.flexc__.w12.displayflex.minHauto.z0.mb7.mobile_px0.mobile_pr0.mobile_pl0.mobile_w12 > div.flexc__.mt3.bgwhite.displayflex.border1.borderColorwhite-8.shadow1.positionrelative.grow1 > div.w12.displayflex.positionrelative.grow1.ipadpro_pl0.ipadpro_pr0.desktop_flexc__ > div > div.bgwhite.flexcb_.mt0.displayflex.grow1.bl1_solid.borderColorgray.desktop_ml0.desktop_mb0.desktop_bordernone > div.pt4.pl3.grow1.flexr_bf"
	var buf1, buf2 []byte
	if err := chromedp.Run(ctx2, func() chromedp.Tasks {
		return chromedp.Tasks{
			chromedp.Emulate(device.KindleFireHDXlandscape),
			chromedp.Navigate(linkURL),
			chromedp.WaitReady("body > #root "),
			chromedp.Sleep(4 * time.Second),
			chromedp.SetAttributeValue(selHeader, "style", "display:none"),
			chromedp.SetAttributeValue(sel1, "style", "margin: 25px 14px"),
			chromedp.ActionFunc(screenshotWithoutPopups2(sel1, &buf1)),
			chromedp.SetAttributeValue(sel2, "style", "margin: 0"),
			chromedp.ActionFunc(screenshotWithoutPopups2(sel2, &buf2)),
		}
	}()); err != nil {
		log.Println(err)
		return nil
	}
	var src image.Image
	if err := glueForTipranks2(buf1, buf2, &src); err != nil {
		log.Println(err)
		return nil
	}
	buf1, buf2 = nil, nil
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

func screenshotWithoutPopups2(sel string, buf *[]byte) func(context.Context) error {
	var fn func(context.Context) error
	fn = func(ctx context.Context) error {
		if err := chromedp.Screenshot(sel, buf, chromedp.NodeVisible).Do(ctx); err != nil {
			return err
		}
		var isPopup1, isPopup2, isOverlay bool
		if err := hidePopup2(ctx, "body #gtm_popup_blocker_iframe", &isPopup1); err != nil {
			return err
		}
		if err := hidePopup2(ctx, "body > #popup-ios-modal-v4", &isPopup2); err != nil {
			return err
		}
		if err := hideOverlay2(ctx, &isOverlay); err != nil {
			return err
		}
		if isPopup1 || isPopup2 || isOverlay {
			return fn(ctx)
		}
		return nil
	}
	return fn
}

func hidePopup2(ctx context.Context, sel string, isPopup *bool) error {
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

func hideOverlay2(ctx context.Context, isOverlay *bool) error {
	const classList = `document.querySelector("#tr-stock-page-content").classList`
	var evaluateResult []byte
	if err := chromedp.Evaluate(classList+`.contains("overlay")`, &evaluateResult).Do(ctx); err != nil {
		return err
	}
	if string(evaluateResult) == "false" {
		return nil
	}
	if err := chromedp.Evaluate(classList+`.remove("overlay")`, &evaluateResult).Do(ctx); err != nil {
		return err
	}
	*isOverlay = true
	return nil
}

func glueForTipranks2(buf1, buf2 []byte, src *image.Image) error {
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
