package screenshot

import (
	"bytes"
	"context"
	"image"
	"image/draw"
	"image/png"
	"log"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
)

func init() {
	// image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
}

// MakeScreenshotForMarketBeat description
func MakeScreenshotForMarketBeat(linkURL string) []byte {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	var buf1, buf2 []byte
	if err := chromedp.Run(ctx, makeScreenshotForMarketBeat(linkURL, &buf1, &buf2)); err != nil {
		log.Fatal(err)
	}
	img1, _, err := image.Decode(bytes.NewReader(buf1))
	if err != nil {
		log.Fatal(err)
	}
	img2, _, err := image.Decode(bytes.NewReader(buf2))
	if err != nil {
		log.Fatal(err)
	}
	//starting position of the second image (bottom left)
	sp2 := image.Point{0, img1.Bounds().Dy()}
	//new rectangle for the second image
	r2 := image.Rectangle{sp2, sp2.Add(img2.Bounds().Size())}
	//rectangle for the big image
	r1 := image.Rectangle{image.Point{0, 0}, r2.Max}
	src := image.NewRGBA(r1)
	draw.Draw(src, img1.Bounds(), img1, image.Point{0, 0}, draw.Src)
	draw.Draw(src, r2, img2, image.Point{0, 0}, draw.Src)

	// // new size of image
	// dr := image.Rect(0, 0, src.Bounds().Max.X/2, src.Bounds().Max.Y/2)
	// // perform resizing
	// res := scaleTo(src, dr, draw.BiLinear)

	// encode
	out := &bytes.Buffer{}
	if err := png.Encode(out, src); err != nil {
		log.Fatal(err)
	}
	// var opt jpeg.Options
	// opt.Quality = 85
	// jpeg.Encode(out, rgba, &opt)
	return out.Bytes()
}

func makeScreenshotForMarketBeat(linkURL string, res1, res2 *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Emulate(device.KindleFireHDX),
		chromedp.Navigate(linkURL),
		chromedp.WaitReady(`body > div > footer`),
		chromedp.WaitVisible("#optinform-modal a"),
		chromedp.Click("#optinform-modal a", chromedp.NodeVisible),

		chromedp.Click("#liInsiderTrades > a", chromedp.NodeVisible),
		chromedp.WaitReady(`body > div > footer`),
		chromedp.WaitVisible("#insiderChart"),
		chromedp.SetAttributeValue("#insiderChart > #svg > #yTextGroup > g.footnote", "style", "display:none"),
		chromedp.Screenshot("#insiderChart", res1, chromedp.NodeVisible),

		chromedp.Click("#liInstutionalOwnership > a", chromedp.NodeVisible),
		chromedp.WaitReady(`body > div > footer`),
		chromedp.WaitVisible("#SECChart"),
		chromedp.SetAttributeValue("#SECChart > #svg > #yTextGroup > g.footnote", "style", "display:none"),
		chromedp.Screenshot("#SECChart", res2, chromedp.NodeVisible),
	}
}

//
// for RGBA images
//

// src   - source image
// rect  - size we want
// scale - scaler
// func scaleTo(src image.Image,
// 	rect image.Rectangle, scale draw.Scaler) image.Image {
// 	dst := image.NewRGBA(rect)
// 	scale.Scale(dst, rect, src, src.Bounds(), draw.Over, nil)
// 	return dst
// }
