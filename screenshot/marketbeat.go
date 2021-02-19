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
	r := image.Rectangle{image.Point{0, 0}, r2.Max}
	rgba := image.NewRGBA(r)
	draw.Draw(rgba, img1.Bounds(), img1, image.Point{0, 0}, draw.Src)
	draw.Draw(rgba, r2, img2, image.Point{0, 0}, draw.Src)
	// encode
	out := &bytes.Buffer{}
	if err := png.Encode(out, rgba); err != nil {
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
