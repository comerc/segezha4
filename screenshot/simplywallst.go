package screenshot

// document.querySelector("#root").style.filter = "none"
// document.querySelectorAll("[data-cy-id='modal-ModalPortal-0-']")[0].style.zIndex = -1

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

// MakeScreenshotForSimplyWallSt description
func MakeScreenshotForSimplyWallSt(linkURL string) ([]byte, []byte) {
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
		return nil, nil
	}
	const average = 14
	ctx2, cancel2 := context.WithTimeout(ctx1, utils.GetTimeout(average))
	defer cancel2()
	// selHeader := "[data-cy-id='header-cover-image']"
	// selNav := "#root > div > nav"
	// script := `
	// (css) => {
	// 	const style = document.createElement('style');
	// 	style.type = 'text/css';
	// 	style.appendChild(document.createTextNode(css));
	// 	document.head.appendChild(style);
	// 	return true;
	// }
	// `
	// var buf1, buf2, buf3, buf4 []byte
	var out1, out2 []byte
	sel1 := "[data-cy-id='overview-summary-snowflake'] > div"
	sel2 := "[data-cy-id='report-sub-section-share-price-vs-fair-value'] > div > div > svg"
	if err := chromedp.Run(ctx2, func() chromedp.Tasks {
		return chromedp.Tasks{
			network.SetBlockedURLS([]string{
				"https://*.stream-io-api.com/*",
				"https://www.google-analytics.com/analytics.js",
				"https://*.cloudfront.net/*",
				"https://static.hotjar.com/*",
				"https://embed.typeform.com/embed.js",
				"https://*.statuspage.io/embed/frame",
				"https://sentry.io/*",
				"https://graphql.simplywall.st/graphql",
			}),
			chromedp.Emulate(device.KindleFireHDX),
			chromedp.Navigate(linkURL),
			chromedp.WaitReady("body"),
			// chromedp.SetAttributeValue("#root", "style", "margin: 0 -16px"),
			// chromedp.SetAttributeValue("#root > div > div", "style", "display:none"),
			// chromedp.SetAttributeValue("#root > div", "style", "padding-top: env(safe-area-inset-top)"),
			// chromedp.SetAttributeValue(selHeader, "style", "display:none"),
			// chromedp.SetAttributeValue(selNav, "style", "display:none"),
			// chromedp.ActionFunc(hideIfExists("iframe")),
			// chromedp.PollFunction(script, nil, chromedp.WithPollingArgs("#root h3:before { display:none }")),
			// chromedp.ActionFunc(AddCSS),
			chromedp.Screenshot(sel1, &out1, chromedp.NodeVisible),
			chromedp.SetAttributeValue("#company-report > div", "style", "display:none"),
			chromedp.SetAttributeValue(sel2, "style", "margin: 20px 0"),
			chromedp.Screenshot(sel2, &out2, chromedp.NodeVisible),
		}
	}()); err != nil {
		log.Println(err)
		return nil, nil
	}
	return out1, out2
	// // !!!!
	// var nodes []*cdp.Node
	// if err := chromedp.Run(ctx2, func() chromedp.Tasks {
	// 	return chromedp.Tasks{
	// 		chromedp.Nodes("section[data-variant]", &nodes, chromedp.AtLeast(0)),
	// 	}
	// }()); err != nil {
	// 	log.Println(err)
	// } else {
	// 	if len(nodes) == 0 {
	// 		if err := takeScreenshotForSimplyWallSt(ctx2, &buf1, &buf2, &buf3, &buf4); err != nil {
	// 			log.Println(err)
	// 		}
	// 	} else {
	// 		if nodes[0].Attributes[0] == "data-variant" && nodes[0].Attributes[1] == "c" {
	// 			if err := takeScreenshotForSimplyWallStVariantC(ctx2, &buf1, &buf2, &buf3, &buf4); err != nil {
	// 				log.Println(err)
	// 			}
	// 		} else {
	// 			if err := takeScreenshotForSimplyWallStOtherVariant(ctx2, &buf1, &buf2, &buf3, &buf4); err != nil {
	// 				log.Println(err)
	// 			}
	// 		}
	// 	}
	// }
	// // !!!!
	// var out1, out2 []byte
	// if len(buf2) == 0 {
	// 	out1 = buf1
	// } else {
	// 	out1 = getOut(buf1, buf2)
	// }
	// out2 = getOut(buf3, buf4)
	// return out1, out2
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

func takeScreenshotForSimplyWallSt(ctx context.Context, buf1, buf2, buf3, buf4 *[]byte) error {
	selSnowflake := "[data-cy-id='company-summary-snowflake']"
	selFairValue := "[data-cy-id='report-sub-section-share-price-vs-fair-value']"
	selERGrowth := "[data-cy-id='report-sub-section-earnings-and-revenue-growth-forecasts']"
	selFutureGrowth := "[data-cy-id='report-sub-section-analyst-future-growth-forecasts']"
	if err := chromedp.Run(ctx, func() chromedp.Tasks {
		return chromedp.Tasks{
			chromedp.SetAttributeValue(selSnowflake+" > h4", "style", "display:none"),
			chromedp.SetAttributeValue(selSnowflake+" > p", "style", "display:none"),
			chromedp.SetAttributeValue(selSnowflake+" > div:nth-child(4)", "style", "display:none"),
			chromedp.SetAttributeValue(selSnowflake, "style", "margin:0"),
			chromedp.Screenshot(selSnowflake, buf1, chromedp.NodeVisible),
			chromedp.SetAttributeValue("#root article > div:nth-child(1)", "style", "display:none"),
			chromedp.SetAttributeValue("#root article > div:nth-child(2)", "style", "display:none"),
		}
	}()); err != nil {
		return err
	}
	var nodes []*cdp.Node
	if err := chromedp.Run(ctx, func() chromedp.Tasks {
		return chromedp.Tasks{
			chromedp.Nodes(selFairValue+" > div > div > div", &nodes, chromedp.AtLeast(0)),
		}
	}()); err != nil {
		return err
	}
	if len(nodes) != 0 {
		if err := chromedp.Run(ctx, func() chromedp.Tasks {
			return chromedp.Tasks{
				chromedp.SetAttributeValue(selFairValue, "style", "padding: 8px"),
				chromedp.SetAttributeValue(selFairValue+" > div > div > div:nth-child(2)", "style", "display: none"),
				chromedp.SetAttributeValue(selFairValue+" > div > div:nth-child(2)", "style", "display: none"),
				chromedp.Screenshot(selFairValue, buf2, chromedp.NodeVisible),
			}
		}()); err != nil {
			return err
		}
	}
	if err := chromedp.Run(ctx, func() chromedp.Tasks {
		return chromedp.Tasks{
			chromedp.SetAttributeValue("#root article > div:nth-child(4)", "style", "display:none"),
			chromedp.SetAttributeValue(selERGrowth, "style", "padding: 8px"),
			chromedp.SetAttributeValue(selERGrowth+" > div > div > div > div > div > div", "style", "display: none"),
			chromedp.SetAttributeValue(selERGrowth+" > div > div > div:nth-child(1) > div:nth-child(2)", "style", "display: none"),
			chromedp.SetAttributeValue(selERGrowth+" > div > div > div:nth-child(2)", "style", "display: none"),
			chromedp.Screenshot(selERGrowth, buf3, chromedp.NodeVisible),
			chromedp.SetAttributeValue(selFutureGrowth, "style", "padding: 8px"),
			chromedp.SetAttributeValue(selFutureGrowth+" > h3", "style", "display: none"),
			chromedp.SetAttributeValue(selFutureGrowth+" > div > div > div:nth-child(2)", "style", "display: none"),
			chromedp.SetAttributeValue(selFutureGrowth+" > div:nth-child(2) > div:nth-child(2)", "style", "display: none"),
			chromedp.Screenshot(selFutureGrowth, buf4, chromedp.NodeVisible),
		}
	}()); err != nil {
		return err
	}
	return nil
}

func takeScreenshotForSimplyWallStOtherVariant(ctx context.Context, buf1, buf2, buf3, buf4 *[]byte) error {
	selCompany := "[data-cy-id='company-summary'] > div"
	// selFairValue := "[data-cy-id='report-sub-section-share-price-vs-fair-value']"
	// selERGrowth := "[data-cy-id='report-sub-section-earnings-and-revenue-growth-forecasts']"
	// selFutureGrowth := "[data-cy-id='report-sub-section-analyst-future-growth-forecasts']"
	if err := chromedp.Run(ctx, func() chromedp.Tasks {
		return chromedp.Tasks{
			chromedp.SetAttributeValue(selCompany, "style", "padding: 0; min-height: 0;"),
			chromedp.SetAttributeValue(selCompany+" > div > h2", "style", "display: none"),
			chromedp.SetAttributeValue(selCompany+" > div > [data-cy-id='company-summary-desc']", "style", "display: none"),
			chromedp.SetAttributeValue(selCompany+" > div > [data-cy-id='company-about-button']", "style", "display: none"),
			chromedp.SetAttributeValue(selCompany+" > div > [data-cy-id='risk-reward-wrapper']", "style", "display: none"),
			chromedp.SetAttributeValue(selCompany+" > div > div > div", "style", "height: 280px"),
			chromedp.SetAttributeValue(selCompany+" > div > div > div > div", "style", "max-width: unset; width: 280px; margin: 0 auto;"),
			chromedp.SetAttributeValue(selCompany+" > div > div > div > div > div:nth-child(2)", "style", "display: none"),
			chromedp.SetAttributeValue(selCompany+" > div > div > div:nth-child(3)", "style", "display: none"),
			chromedp.SetAttributeValue(selCompany+" > div > div", "style", "max-width: unset; margin: 0; float: unset; flex: unset;"),
			// chromedp.SetAttributeValue("[data-cy-id='summary-risk-modal-button']", "style", "display:none"),
			// chromedp.SetAttributeValue("[data-cy-id='overview-summary-snowflake'] > div:nth-child(2)", "style", "display:none"),
			// chromedp.SetAttributeValue("[data-cy-id='chart-action-toggle-data-overview-summary-snowflake']", "style", "display:none"),
			// chromedp.SetAttributeValue("[data-cy-id='chart-action-help-overview-summary-snowflake']", "style", "display:none"),
			// chromedp.SetAttributeValue("[data-cy-id='chart-action-more-overview-summary-snowflake']", "style", "display:none"),
			chromedp.Screenshot(selCompany, buf1, chromedp.NodeVisible),
		}
	}()); err != nil {
		return err
	}
	// {
	// 	var nodes []*cdp.Node
	// 	if err := chromedp.Run(ctx, func() chromedp.Tasks {
	// 		return chromedp.Tasks{
	// 			chromedp.Nodes("#root article > section", &nodes, chromedp.AtLeast(0)),
	// 		}
	// 	}()); err != nil {
	// 		return err
	// 	}
	// 	log.Print(len(nodes))
	// 	if err := chromedp.Run(ctx, func() chromedp.Tasks {
	// 		return chromedp.Tasks{
	// 			chromedp.ActionFunc(func(ctx context.Context) error {
	// 				for i, node := range nodes {
	// 					if i == 7 {
	// 						break
	// 					}
	// 					log.Printf("%#v", node.Attributes)
	// 					if err := dom.SetAttributeValue(nodes[i].NodeID, "style", "display:none").Do(ctx); err != nil {
	// 						return err
	// 					}
	// 				}
	// 				return nil
	// 			}),
	// 		}
	// 	}()); err != nil {
	// 		return err
	// 	}
	// }
	// var nodes []*cdp.Node
	// if err := chromedp.Run(ctx, func() chromedp.Tasks {
	// 	return chromedp.Tasks{
	// 		chromedp.Nodes(selFairValue+" > div > div > div", &nodes, chromedp.AtLeast(0)),
	// 	}
	// }()); err != nil {
	// 	return err
	// }
	// log.Print(len(nodes))
	// if len(nodes) != 0 {
	// 	if err := chromedp.Run(ctx, func() chromedp.Tasks {
	// 		return chromedp.Tasks{
	// 			chromedp.SetAttributeValue(selFairValue, "style", "padding: 8px"),
	// 			chromedp.SetAttributeValue(selFairValue+" > div > div > div:nth-child(2)", "style", "display: none"),
	// 			chromedp.SetAttributeValue(selFairValue+" > div > div:nth-child(2)", "style", "display: none"),
	// 			chromedp.Screenshot(selFairValue, buf2, chromedp.NodeVisible),
	// 		}
	// 	}()); err != nil {
	// 		return err
	// 	}
	// }
	// if err := chromedp.Run(ctx, func() chromedp.Tasks {
	// 	return chromedp.Tasks{
	// 		chromedp.SetAttributeValue("#root article > section:nth-child(5)", "style", "display:none"),
	// 		chromedp.SetAttributeValue(selERGrowth, "style", "padding: 8px"),
	// 		chromedp.SetAttributeValue(selERGrowth+" > div > div > div > div > div > div", "style", "display: none"),
	// 		chromedp.SetAttributeValue(selERGrowth+" > div > div > div:nth-child(1) > div:nth-child(2)", "style", "display: none"),
	// 		chromedp.SetAttributeValue(selERGrowth+" > div > div > div:nth-child(2)", "style", "display: none"),
	// 		chromedp.Screenshot(selERGrowth, buf3, chromedp.NodeVisible),
	// 		// chromedp.SetAttributeValue(selERGrowth, "style", "display:none"),
	// 		chromedp.SetAttributeValue(selFutureGrowth, "style", "padding: 8px"),
	// 		chromedp.SetAttributeValue(selFutureGrowth+" > h3", "style", "display: none"),
	// 		chromedp.SetAttributeValue(selFutureGrowth+" > div > div > div:nth-child(2)", "style", "display: none"),
	// 		chromedp.SetAttributeValue(selFutureGrowth+" > div:nth-child(2) > div:nth-child(2)", "style", "display: none"),
	// 		chromedp.Screenshot(selFutureGrowth, buf4, chromedp.NodeVisible),
	// 	}
	// }()); err != nil {
	// 	return err
	// }
	return nil
}

func takeScreenshotForSimplyWallStVariantC(ctx context.Context, buf1, buf2, buf3, buf4 *[]byte) error {
	selCompany := "[data-cy-id='company-summary'] > div"
	selFairValue := "[data-cy-id='report-sub-section-share-price-vs-fair-value']"
	selERGrowth := "[data-cy-id='report-sub-section-earnings-and-revenue-growth-forecasts']"
	selFutureGrowth := "[data-cy-id='report-sub-section-analyst-future-growth-forecasts']"
	if err := chromedp.Run(ctx, func() chromedp.Tasks {
		return chromedp.Tasks{
			chromedp.SetAttributeValue(selCompany, "style", "padding: 0; min-height: 0;"),
			chromedp.SetAttributeValue(selCompany+" > div > h2", "style", "display: none"),
			chromedp.SetAttributeValue(selCompany+" > div > [data-cy-id='company-summary-desc']", "style", "display: none"),
			chromedp.SetAttributeValue(selCompany+" > div > [data-cy-id='risk-reward-wrapper']", "style", "display: none"),
			chromedp.SetAttributeValue(selCompany+" > div > div > div", "style", "height: 280px"),
			chromedp.SetAttributeValue(selCompany+" > div > div > div > div", "style", "max-width: unset; width: 280px; margin: 0 auto;"),
			chromedp.SetAttributeValue(selCompany+" > div > div > div > div > div:nth-child(2)", "style", "display: none"),
			chromedp.SetAttributeValue(selCompany+" > div > div > div:nth-child(3)", "style", "display: none"),
			chromedp.SetAttributeValue(selCompany+" > div > div", "style", "max-width: unset; margin: 0; float: unset; flex: unset;"),
			// chromedp.SetAttributeValue("[data-cy-id='summary-risk-modal-button']", "style", "display:none"),
			// chromedp.SetAttributeValue("[data-cy-id='overview-summary-snowflake'] > div:nth-child(2)", "style", "display:none"),
			// chromedp.SetAttributeValue("[data-cy-id='chart-action-toggle-data-overview-summary-snowflake']", "style", "display:none"),
			// chromedp.SetAttributeValue("[data-cy-id='chart-action-help-overview-summary-snowflake']", "style", "display:none"),
			// chromedp.SetAttributeValue("[data-cy-id='chart-action-more-overview-summary-snowflake']", "style", "display:none"),
			chromedp.Screenshot(selCompany, buf1, chromedp.NodeVisible),
			chromedp.SetAttributeValue("#root article > section:nth-child(1)", "style", "display:none"),
			chromedp.SetAttributeValue("#root article > section:nth-child(2)", "style", "display:none"),
			chromedp.SetAttributeValue("#root article > section:nth-child(3)", "style", "display:none"),
			chromedp.SetAttributeValue("#root article > section.dabTGn", "style", "display:none"),
		}
	}()); err != nil {
		return err
	}
	var nodes []*cdp.Node
	if err := chromedp.Run(ctx, func() chromedp.Tasks {
		return chromedp.Tasks{
			chromedp.Nodes(selFairValue+" > div > div > div", &nodes, chromedp.AtLeast(0)),
		}
	}()); err != nil {
		return err
	}
	if len(nodes) != 0 {
		if err := chromedp.Run(ctx, func() chromedp.Tasks {
			return chromedp.Tasks{
				chromedp.SetAttributeValue(selFairValue, "style", "padding: 8px"),
				chromedp.SetAttributeValue(selFairValue+" > div > div > div:nth-child(2)", "style", "display: none"),
				chromedp.SetAttributeValue(selFairValue+" > div > div:nth-child(2)", "style", "display: none"),
				chromedp.Screenshot(selFairValue, buf2, chromedp.NodeVisible),
			}
		}()); err != nil {
			return err
		}
	}
	if err := chromedp.Run(ctx, func() chromedp.Tasks {
		return chromedp.Tasks{
			chromedp.SetAttributeValue("#root article > section:nth-child(5)", "style", "display:none"),
			chromedp.SetAttributeValue(selERGrowth, "style", "padding: 8px"),
			chromedp.SetAttributeValue(selERGrowth+" > div > div > div > div > div > div", "style", "display: none"),
			chromedp.SetAttributeValue(selERGrowth+" > div > div > div:nth-child(1) > div:nth-child(2)", "style", "display: none"),
			chromedp.SetAttributeValue(selERGrowth+" > div > div > div:nth-child(2)", "style", "display: none"),
			chromedp.Screenshot(selERGrowth, buf3, chromedp.NodeVisible),
			// chromedp.SetAttributeValue(selERGrowth, "style", "display:none"),
			chromedp.SetAttributeValue(selFutureGrowth, "style", "padding: 8px"),
			chromedp.SetAttributeValue(selFutureGrowth+" > h3", "style", "display: none"),
			chromedp.SetAttributeValue(selFutureGrowth+" > div > div > div:nth-child(2)", "style", "display: none"),
			chromedp.SetAttributeValue(selFutureGrowth+" > div:nth-child(2) > div:nth-child(2)", "style", "display: none"),
			chromedp.Screenshot(selFutureGrowth, buf4, chromedp.NodeVisible),
		}
	}()); err != nil {
		return err
	}
	return nil
}
