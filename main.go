package main

import (
	"image"
	"image/gif"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/nfnt/resize"
	"github.com/oliamb/cutter"
)

// TODO resolve is bug
//  2021/01/12 17:24:04 fail to encode gif lzw: input byte too large for the litWidth

func main() {
	g := GetImage("sample/ethereum.gif")
	//spew.Dump(len(g.Image), g.BackgroundIndex, g.Config, g.Delay, g.Disposal, g.LoopCount)

	var h int
	var w int

	for i, img := range g.Image {
		// TODO Resizer changes the image type from image.Paletted to image.RGBA64
		ri := ResizeSquare(80, img)
		ci := Crop(ri)

		pci := image.Paletted{
			Pix:     ci.(*image.RGBA64).Pix,
			Stride:  ci.(*image.RGBA64).Stride,
			Rect:    ci.(*image.RGBA64).Rect,
			Palette: img.Palette,
		}

		g.Image[i] = &pci

		if h == 0 {
			h = pci.Bounds().Dy()
			spew.Dump(pci.Rect, ri.Bounds())
		}

		if w == 0 {
			w = pci.Bounds().Dx()
		}
	}

	g.Config.Height = h
	g.Config.Width = w

	spew.Dump(g.Config.Width, g.Config.Height, h, w)

	RenderImage(g, "sample/out.gif")
}

func GetImage(fileName string) *gif.GIF {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}

	img, err := gif.DecodeAll(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	return img
}

func ResizeSquare(size uint, img image.Image) image.Image {
	return resize.Resize(size, size, img, resize.NearestNeighbor)
}

func Crop(img image.Image) image.Image {
	var sl int
	if img.Bounds().Max.X > img.Bounds().Max.Y {
		sl = img.Bounds().Max.Y
	} else {
		sl = img.Bounds().Max.X
	}

	croppedImg, err := cutter.Crop(img, cutter.Config{
		Width:  sl,
		Height: sl,
		Mode:   cutter.Centered,
	})
	if err != nil {
		log.Fatal(err)
	}

	return croppedImg
}

func RenderImage(g *gif.GIF, filename string) {
	out, err := os.Create(filename)
	if err != nil {
		log.Fatal("fail to create new file", err)
	}
	defer out.Close()

	err = gif.EncodeAll(out, g)
	if err != nil {
		log.Fatal("fail to encode gif ", err)
	}

	fi, _ := out.Stat()
	spew.Dump(fi)
}