package screenshot

// document.querySelector("#root").style.filter = "none"
// document.querySelectorAll("[data-cy-id='modal-ModalPortal-0-']")[0].style.zIndex = -1

import (
	"bytes"
	"context"
	"image"
	"image/png"
	"log"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
	"github.com/comerc/segezha4/utils"
)

// MakeScreenshotForSimplyWallSt description
func MakeScreenshotForSimplyWallSt(linkURL string) ([]byte, []byte) {
	ctx1, cancel1 := chromedp.NewContext(context.Background())
	defer cancel1()
	// start the browser without a timeout
	if err := chromedp.Run(ctx1); err != nil {
		log.Println(err)
		return nil, nil
	}
	const average = 14
	ctx2, cancel2 := context.WithTimeout(ctx1, utils.GetTimeout(average))
	defer cancel2()
	selNav := "#root > div > nav"
	selSnowflake := "[data-cy-id='company-summary-snowflake']"
	selFairValue := "[data-cy-id='report-sub-section-share-price-vs-fair-value']"
	selERGrowth := "[data-cy-id='report-sub-section-earnings-and-revenue-growth-forecasts']"
	selFutureGrowth := "[data-cy-id='report-sub-section-analyst-future-growth-forecasts']"
	// script := `
	// (css) => {
	// 	const style = document.createElement('style');
	// 	style.type = 'text/css';
	// 	style.appendChild(document.createTextNode(css));
	// 	document.head.appendChild(style);
	// 	return true;
	// }
	// `
	var buf1, buf2, buf3, buf4 []byte
	if err := chromedp.Run(ctx2, func() chromedp.Tasks {
		return chromedp.Tasks{
			chromedp.Emulate(device.KindleFireHDX),
			chromedp.Navigate(linkURL),
			chromedp.WaitReady("body"),
			chromedp.SetAttributeValue("#root", "style", "margin: -16px"),
			chromedp.SetAttributeValue(selNav, "style", "display:none"),
			// chromedp.PollFunction(script, nil, chromedp.WithPollingArgs("#root h3:before { display:none }")),
			chromedp.ActionFunc(AddCSS),

			chromedp.SetAttributeValue(selSnowflake+" > h4", "style", "display:none"),
			chromedp.SetAttributeValue(selSnowflake+" > p", "style", "display:none"),
			chromedp.SetAttributeValue(selSnowflake+" > div:nth-child(4)", "style", "display:none"),
			chromedp.SetAttributeValue(selSnowflake, "style", "margin:0"),
			chromedp.Screenshot(selSnowflake, &buf1, chromedp.NodeVisible),

			chromedp.SetAttributeValue("#root article > div:nth-child(1)", "style", "display:none"),
			chromedp.SetAttributeValue("#root article > div:nth-child(2)", "style", "display:none"),

			chromedp.SetAttributeValue(selFairValue, "style", "padding: 8px"),
			chromedp.SetAttributeValue(selFairValue+" > div > div > div:nth-child(2)", "style", "display: none"),
			chromedp.SetAttributeValue(selFairValue+" > div > div:nth-child(2)", "style", "display: none"),
			chromedp.Screenshot(selFairValue, &buf2, chromedp.NodeVisible),

			chromedp.SetAttributeValue("#root article > div:nth-child(4)", "style", "display:none"),

			chromedp.SetAttributeValue(selERGrowth, "style", "padding: 8px"),
			chromedp.SetAttributeValue(selERGrowth+" > div > div > div > div > div > div", "style", "display: none"),
			chromedp.SetAttributeValue(selERGrowth+" > div > div > div:nth-child(1) > div:nth-child(2)", "style", "display: none"),
			chromedp.SetAttributeValue(selERGrowth+" > div > div > div:nth-child(2)", "style", "display: none"),
			chromedp.Screenshot(selERGrowth, &buf3, chromedp.NodeVisible),

			chromedp.SetAttributeValue(selFutureGrowth, "style", "padding: 8px"),
			chromedp.SetAttributeValue(selFutureGrowth+" > h3", "style", "display: none"),
			chromedp.SetAttributeValue(selFutureGrowth+" > div > div > div:nth-child(2)", "style", "display: none"),
			chromedp.SetAttributeValue(selFutureGrowth+" > div:nth-child(2) > div:nth-child(2)", "style", "display: none"),
			chromedp.Screenshot(selFutureGrowth, &buf4, chromedp.NodeVisible),
		}
	}()); err != nil {
		log.Println(err)
		return nil, nil
	}

	return getOut(buf1, buf2), getOut(buf3, buf4)
}

func getOut(buf1, buf2 []byte) []byte {
	var src image.Image
	if err := glueForSimplyWallSt(buf1, buf2, &src); err != nil {
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

func glueForSimplyWallSt(buf1, buf2 []byte, src *image.Image) error {
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
	img1 = nil
	img2 = nil
	return nil
}

func AddCSS(ctx context.Context) error {
	css := "#root h3:before { display:none }"
	script := `
	(() => {
		const style = document.createElement('style');
		style.type = 'text/css';
		style.appendChild(document.createTextNode(` + "`" + css + "`" + `));
		document.head.appendChild(style);
	})()
	`
	var evaluateResult []byte
	err := chromedp.Evaluate(script, &evaluateResult).Do(ctx)
	return err
}
