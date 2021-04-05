package screenshot

import (
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"net/http"
	"os"

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

func getWebSocketDebuggerUrl() string {
	headlessIP := os.Getenv("HEADLESS_IP")
	if headlessIP == "" {
		headlessIP = "localhost"
	}
	fmt.Println(headlessIP)
	resp, err := http.Get(fmt.Sprintf("http://%s:9222/json/version", headlessIP))

	// resp, err := http.Get("http://172.16.0.42:9222/json/version")
	if err != nil {
		log.Fatal(err)
	}

	var result map[string]interface{}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatal(err)
	}
	return result["webSocketDebuggerUrl"].(string)
}

// func GetChrome(timeout time.Duration) (context.Context, context.CancelFunc) {
// 	// o := append(chromedp.DefaultExecAllocatorOptions[:],
// 	// 	// chromedp.ProxyServer("socks5://138.59.207.118:9076"),
// 	// 	chromedp.Flag("blink-settings", "imagesEnabled=false"),
// 	// )
// 	// ctx, cancel := chromedp.NewExecAllocator(context.Background(), o...)
// 	// defer cancel()
// 	// ctx1, cancel1 := chromedp.NewContext(ctx)
// 	// defer cancel1()
// 	// // ctx1, cancel1 := chromedp.NewContext(context.Background())
// 	// // defer cancel1()
// 	// // start the browser without a timeout
// 	// if err := chromedp.Run(ctx1); err != nil {
// 	// 	log.Println(err)
// 	// 	return nil
// 	// }
// 	// ctx2, cancel2 := context.WithTimeout(ctx1, 100*time.Second)
// 	// defer cancel2()

// 	// ctx, cancel := chromedp.NewRemoteAllocator(context.Background(), getWebSocketDebuggerUrl())
// 	// defer cancel()
// 	// ctx1, cancel1 := chromedp.NewContext(ctx)
// 	// defer cancel1()
// 	// // start the browser without a timeout
// 	// if err := chromedp.Run(ctx1); err != nil {
// 	// 	log.Println(err)
// 	// 	// return nil, nil
// 	// }
// 	// ctx2, cancel2 := context.WithTimeout(ctx1, timeout)
// 	// defer cancel2()
// }
