package main

import (
	"bytes"
	"context"
	"github.com/disintegration/imaging"
	"image"
	"image/color"
	"io/ioutil"
	"log"

	"github.com/chromedp/chromedp"
)

func main() {
	// create context
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		// chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	// capture screenshot of an element
	var buf []byte
	if err := chromedp.Run(ctx, elementScreenshot(`https://pkg.go.dev/`, &buf)); err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile("elementScreenshot.png", buf, 0o644); err != nil {
		log.Fatal(err)
	}

	// capture entire browser viewport, returning png with quality=90
	if err := chromedp.Run(ctx, fullScreenshot(`https://google.ca/`, 80, &buf)); err != nil {
		log.Fatal(err)
	}
	src, _ , err := image.Decode(bytes.NewReader(buf))
	if err != nil {
		log.Fatal(err)
	}
	src = imaging.Resize(src, 300, 0, imaging.Lanczos)
	src = imaging.CropAnchor(src, 300, 300, imaging.Top)
	dst := imaging.New(300, 300, color.YCbCr{Y: 4, Cb: 2, Cr: 2})
	dst = imaging.Paste(dst, src, image.Pt(0,0))
	err = imaging.Save(dst, "fittedImage.jpg")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("wrote elementScreenshot.png and fullScreenshot.png")
}

// elementScreenshot takes a screenshot of a specific element.
func elementScreenshot(urlstr string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.EmulateViewport(300,300),
		chromedp.Navigate(urlstr),
		chromedp.CaptureScreenshot(res),
	}
}

// fullScreenshot takes a screenshot of the entire browser viewport.
//
// Note: chromedp.FullScreenshot overrides the device's emulation settings. Reset
func fullScreenshot(urlstr string, quality int, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.FullScreenshot(res, quality),
	}
}