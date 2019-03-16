package main

import (
	"fmt"
	"image"
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
	source *ymage
	target *ymage
}

type ymage struct {
	pixels []lang.Color
	width  int
	height int
}

func loadSurface(reader io.Reader) (*surface, error) {
	source, err := loadImage(reader)
	if err != nil {
		return nil, err
	}
	target := &ymage{
		pixels: make([]lang.Color, len(source.pixels)),
		width:  source.width,
		height: source.height,
	}
	return &surface{source, target}, nil
}

func (surf *surface) GetPixel(x int, y int) lang.Color {
	return surf.source.pixels[y*surf.source.width+x]
}

func (surf *surface) SetPixel(x int, y int, col lang.Color) {
	surf.target.pixels[y*surf.target.width+x] = col
}

func (surf *surface) SourceWidth() int {
	return surf.source.width
}

func (surf *surface) SourceHeight() int {
	return surf.source.height
}

func (surf *surface) TargetWidth() int {
	return surf.target.width
}

func (surf *surface) TargetHeight() int {
	return surf.target.height
}

func (surf *surface) Convolute(x, y, width, height int, kernel []lang.Number) lang.Color {
	kernelSum := lang.Number(0)
	r := lang.Number(0)
	g := lang.Number(0)
	b := lang.Number(0)
	a := lang.Number(255)
	kernelIndex := 0
	w := surf.SourceWidth()
	h := surf.SourceHeight()

	for kernelY := 0; kernelY < width; kernelY++ {
		for kernelX := 0; kernelX < width; kernelX++ {
			sourceY := y - (height / 2) + kernelY
			sourceX := x - (width / 2) + kernelX
			if sourceX >= 0 && sourceX < w && sourceY >= 0 && sourceY < h {
				value := kernel[kernelIndex]
				px := surf.GetPixel(sourceX, sourceY)
				r += value * px.R
				g += value * px.G
				b += value * px.B
				kernelSum += value

				if sourceX == x && sourceY == y {
					a = px.A
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

func (surf *surface) mapChannel(x, y, width, height int, kernel []lang.Number, mapper func(lang.Color) lang.Number) []lang.Number {
	result := make([]lang.Number, len(kernel))
	kernelIndex := 0
	w := surf.SourceWidth()
	h := surf.SourceHeight()

	for kernelY := 0; kernelY < width; kernelY++ {
		for kernelX := 0; kernelX < width; kernelX++ {
			sourceY := y - (width / 2) + kernelY
			sourceX := x - (height / 2) + kernelX
			if sourceX >= 0 && sourceX < w && sourceY >= 0 && sourceY < h {
				px := surf.GetPixel(sourceX, sourceY)
				result[kernelIndex] = kernel[kernelIndex] * mapper(px)
			}
			kernelIndex++
		}
	}
	return result
}

func (surf *surface) MapRed(x, y, radius, width int, kernel []lang.Number) []lang.Number {
	return surf.mapChannel(x, y, radius, width, kernel, func(px lang.Color) lang.Number { return px.R })
}

func (surf *surface) MapGreen(x, y, radius, width int, kernel []lang.Number) []lang.Number {
	return surf.mapChannel(x, y, radius, width, kernel, func(px lang.Color) lang.Number { return px.G })
}

func (surf *surface) MapBlue(x, y, radius, width int, kernel []lang.Number) []lang.Number {
	return surf.mapChannel(x, y, radius, width, kernel, func(px lang.Color) lang.Number { return px.B })
}

func (surf *surface) MapAlpha(x, y, radius, width int, kernel []lang.Number) []lang.Number {
	return surf.mapChannel(x, y, radius, width, kernel, func(px lang.Color) lang.Number { return px.A })
}

func (surf *surface) Blt(x, y, width, height int) {
	if width == surf.SourceWidth() && height == surf.SourceHeight() && surf.target.width == surf.source.width && surf.source.height == surf.target.height {
		copy(surf.target.pixels, surf.source.pixels)
	} else {
		for iy := y; iy < height; iy++ {
			for ix := x; ix < width; ix++ {
				index := iy*surf.target.width + ix
				surf.target.pixels[index] = surf.source.pixels[index]
			}
		}
	}
}

func (surf *surface) ResizeTarget(width, height int) {
	surf.target = &ymage{
		width:  width,
		height: height,
		pixels: make([]lang.Color, width*height),
	}
}

func (surf *surface) Flip() {
	surf.source = surf.target
	surf.target = &ymage{
		width:  surf.source.width,
		height: surf.source.height,
		pixels: append([]lang.Color(nil), surf.source.pixels...),
	}
}

func loadImage(reader io.Reader) (*ymage, error) {
	source, encoding, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	log.Printf("Image decoded as %s", encoding)
	target := image.NewNRGBA(source.Bounds())
	draw.Draw(target, target.Bounds(), source, image.Point{0, 0}, draw.Src)

	byteCount := len(target.Pix)
	pixels := make([]lang.Color, byteCount/4)
	j := 0
	for i := 0; i < byteCount; i += 4 {
		pixels[j] = lang.NewRgba(
			lang.Number(target.Pix[i+0]),
			lang.Number(target.Pix[i+1]),
			lang.Number(target.Pix[i+2]),
			lang.Number(target.Pix[i+3]))
		j++
	}

	return &ymage{
		pixels: pixels,
		width:  target.Rect.Dx(),
		height: target.Rect.Dy(),
	}, nil
}

func saveImage(ymg *ymage, targetPath string) error {
	targetFile, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("error creating file %s: %s", targetPath, err)
	}
	defer targetFile.Close()

	img := image.NewNRGBA(image.Rect(0, 0, ymg.width, ymg.height))
	byteCount := len(img.Pix)
	j := 0
	for i := 0; i < byteCount; i += 4 {
		rgba := ymg.pixels[j].Clamp()
		img.Pix[i+0] = byte(rgba.R)
		img.Pix[i+1] = byte(rgba.G)
		img.Pix[i+2] = byte(rgba.B)
		img.Pix[i+3] = byte(rgba.A)
		j++
	}

	if err = png.Encode(targetFile, img); err != nil {
		return fmt.Errorf("error encoding target image: %s", err)
	}
	return nil
}
