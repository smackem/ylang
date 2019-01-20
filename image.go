package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"log"
	"os"

	"github.com/smackem/ylang/internal/lang"
)

type surface struct {
	source *image.NRGBA
	target *image.NRGBA
}

func loadSurface(sourcePath string) (*surface, error) {
	source, err := loadImage(sourcePath)
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

func (surf *surface) Convolute(x int, y int, radius int, length int, kernel []lang.Number) lang.Color {
	kernelSum := lang.Number(0.0)
	r := lang.Number(0.0)
	g := lang.Number(0.0)
	b := lang.Number(0.0)
	a := lang.Number(255)
	kernelIndex := 0
	width := surf.Width()
	height := surf.Height()

	for kernelY := 0; kernelY < length; kernelY++ {
		for kernelX := 0; kernelX < length; kernelX++ {
			sourceY := y - radius + kernelY
			sourceX := x - radius + kernelX
			if sourceX >= 0 && sourceX < width && sourceY >= 0 && sourceY < height {
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

func (surf *surface) Blt(rect lang.Rect) {
	copy(surf.target.Pix, surf.source.Pix)
}

func loadImage(sourcePath string) (*image.NRGBA, error) {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("Could not load %s: %s", sourcePath, err.Error())
	}
	defer sourceFile.Close()

	source, encoding, err := image.Decode(sourceFile)
	if err != nil {
		return nil, fmt.Errorf("Error decoding %s: %s", sourcePath, err.Error())
	}

	log.Printf("Image %s decoded as %s", sourcePath, encoding)

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
