package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io"
	"log"
	"os"

	"github.com/smackem/ylang/internal/lang"
)

type surface struct {
	source *image.NRGBA
	target *image.NRGBA
}

func loadSurface(reader io.Reader) (*surface, error) {
	source, err := loadImage(reader)
	if err != nil {
		return nil, err
	}
	target := &image.NRGBA{
		Pix:    make([]byte, len(source.Pix)),
		Rect:   source.Rect,
		Stride: source.Stride,
	}
	return &surface{source, target}, nil
}

func (surf *surface) GetPixel(x int, y int) lang.Color {
	nrgba := surf.source.NRGBAAt(x, y)
	return lang.NewRgba(lang.Number(nrgba.R), lang.Number(nrgba.G), lang.Number(nrgba.B), lang.Number(nrgba.A))
}

func (surf *surface) SetPixel(x int, y int, col lang.Color) {
	clamped := col.Clamp()
	surf.target.SetNRGBA(x, y, color.NRGBA{
		R: byte(clamped.R),
		G: byte(clamped.G),
		B: byte(clamped.B),
		A: byte(clamped.A),
	})
}

func (surf *surface) Width() int {
	return surf.source.Bounds().Dx()
}

func (surf *surface) Height() int {
	return surf.source.Bounds().Dy()
}

func (surf *surface) Convolute(x int, y int, radius int, width int, kernel []lang.Number) lang.Color {
	kernelSum := lang.Number(0.0)
	r := lang.Number(0.0)
	g := lang.Number(0.0)
	b := lang.Number(0.0)
	a := lang.Number(255)
	kernelIndex := 0
	w := surf.Width()
	h := surf.Height()

	for kernelY := 0; kernelY < width; kernelY++ {
		for kernelX := 0; kernelX < width; kernelX++ {
			sourceY := y - radius + kernelY
			sourceX := x - radius + kernelX
			if sourceX >= 0 && sourceX < w && sourceY >= 0 && sourceY < h {
				value := kernel[kernelIndex]
				px := surf.source.NRGBAAt(sourceX, sourceY)
				r += value * lang.Number(px.R)
				g += value * lang.Number(px.G)
				b += value * lang.Number(px.B)
				kernelSum += value

				if sourceX == x && sourceY == y {
					a = lang.Number(px.A)
				}
			}
			kernelIndex++
		}
	}
	if kernelSum == 0.0 {
		return lang.NewRgba(r, g, b, a)
	}

	return lang.NewRgba(r/kernelSum, g/kernelSum, b/kernelSum, a)
}

func (surf *surface) mapChannel(x, y, radius, width int, kernel []lang.Number, mapper func(color.NRGBA) byte) []lang.Number {
	result := make([]lang.Number, len(kernel))
	kernelIndex := 0
	w := surf.Width()
	h := surf.Height()

	for kernelY := 0; kernelY < width; kernelY++ {
		for kernelX := 0; kernelX < width; kernelX++ {
			sourceY := y - radius + kernelY
			sourceX := x - radius + kernelX
			if sourceX >= 0 && sourceX < w && sourceY >= 0 && sourceY < h {
				px := surf.source.NRGBAAt(sourceX, sourceY)
				result[kernelIndex] = kernel[kernelIndex] * lang.Number(mapper(px))
			}
			kernelIndex++
		}
	}
	return result
}

func (surf *surface) MapRed(x, y, radius, width int, kernel []lang.Number) []lang.Number {
	return surf.mapChannel(x, y, radius, width, kernel, func(px color.NRGBA) byte { return px.R })
}

func (surf *surface) MapGreen(x, y, radius, width int, kernel []lang.Number) []lang.Number {
	return surf.mapChannel(x, y, radius, width, kernel, func(px color.NRGBA) byte { return px.G })
}

func (surf *surface) MapBlue(x, y, radius, width int, kernel []lang.Number) []lang.Number {
	return surf.mapChannel(x, y, radius, width, kernel, func(px color.NRGBA) byte { return px.B })
}

func (surf *surface) MapAlpha(x, y, radius, width int, kernel []lang.Number) []lang.Number {
	return surf.mapChannel(x, y, radius, width, kernel, func(px color.NRGBA) byte { return px.A })
}

func (surf *surface) Blt(rect lang.Rect) {
	if rect.Max.X-rect.Min.X == surf.Width() && rect.Max.Y-rect.Min.Y == surf.Height() {
		copy(surf.source.Pix, surf.target.Pix)
	} else {
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			for x := rect.Min.X; x < rect.Max.X; x++ {
				surf.source.SetNRGBA(x, y, surf.target.NRGBAAt(x, y))
			}
		}
	}
}

func loadImage(reader io.Reader) (*image.NRGBA, error) {
	source, encoding, err := image.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("Error decoding %s: %s", reader, err.Error())
	}

	log.Printf("Image %s decoded as %s", reader, encoding)

	target := image.NewNRGBA(source.Bounds())
	draw.Draw(target, target.Bounds(), source, image.Point{0, 0}, draw.Src)

	return target, nil
}

func saveImage(img image.Image, targetPath string) error {
	targetFile, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("Could not create file %s: %s", targetPath, err)
	}
	defer targetFile.Close()

	if err = png.Encode(targetFile, img); err != nil {
		return fmt.Errorf("Could not encode target image: %s", err)
	}

	return nil
}
