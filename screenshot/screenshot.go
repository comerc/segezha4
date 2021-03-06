package screenshot

import (
	"context"
	"image"
	"image/draw"
	"image/png"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

func init() {
	// image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
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

func hideIfExists(sel string) func(context.Context) error {
	return func(ctx context.Context) error {
		var nodes []*cdp.Node
		if err := chromedp.Nodes(sel, &nodes, chromedp.AtLeast(0)).Do(ctx); err != nil {
			return err
		}
		if len(nodes) == 0 {
			return nil
		}
		if err := chromedp.SetAttributeValue(sel, "style", "display:none").Do(ctx); err != nil {
			return err
		}
		return nil
	}
}
