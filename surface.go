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
	"math"
	"os"

	"github.com/smackem/ylang/internal/lang"
)

type surface struct {
	source        *ymage
	target        *ymage
	sourceHistory []*ymage
	clipRect      image.Rectangle
	log           func(string)
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
	return &surface{
		source:        source,
		target:        target,
		sourceHistory: nil,
	}, nil
}

func (surf *surface) GetPixel(x int, y int) lang.Color {
	return surf.source.pixels[y*surf.source.width+x]
}

func (surf *surface) SetPixel(x int, y int, col lang.Color) {
	if surf.clipRect.Empty() == false {
		pt := image.Point{x, y}
		if pt.In(surf.clipRect) == false {
			return
		}
	}
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

	for kernelY := 0; kernelY < height; kernelY++ {
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

	for kernelY := 0; kernelY < height; kernelY++ {
		for kernelX := 0; kernelX < width; kernelX++ {
			sourceY := y - (height / 2) + kernelY
			sourceX := x - (width / 2) + kernelX
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
	bltRect := image.Rect(x, y, width, height)
	srcRect := image.Rect(0, 0, surf.source.width, surf.source.height)
	trgRect := image.Rect(0, 0, surf.target.width, surf.target.height)

	if bltRect == srcRect && trgRect == srcRect {
		if surf.clipRect.Empty() || surf.clipRect == bltRect {
			copy(surf.target.pixels, surf.source.pixels)
			return
		}
	}

	rect := bltRect.Intersect(srcRect).Intersect(trgRect).Intersect(surf.clipRect)

	for iy := rect.Min.Y; iy < rect.Max.Y; iy++ {
		for ix := rect.Min.X; ix < rect.Max.X; ix++ {
			index := iy*surf.target.width + ix
			surf.target.pixels[index] = surf.source.pixels[index]
		}
	}
}

func (surf *surface) ResizeTarget(width, height int) {
	surf.target = &ymage{
		width:  width,
		height: height,
		pixels: make([]lang.Color, width*height),
	}
	surf.clipRect = image.Rectangle{}
}

func (surf *surface) Flip() int {
	oldSourceID := len(surf.sourceHistory)
	surf.sourceHistory = append(surf.sourceHistory, surf.source)
	surf.source = surf.target
	surf.target = &ymage{
		width:  surf.source.width,
		height: surf.source.height,
		pixels: append([]lang.Color(nil), surf.source.pixels...),
	}
	surf.clipRect = image.Rectangle{}
	return oldSourceID
}

func (surf *surface) Recall(imageID int) error {
	if imageID >= len(surf.sourceHistory) {
		return fmt.Errorf("unknown context %d - cannot recall", imageID)
	}
	surf.source = surf.sourceHistory[imageID]
	return nil
}

func (surf *surface) ClipRect() image.Rectangle {
	return surf.clipRect
}

func (surf *surface) SetClipRect(rect image.Rectangle) {
	surf.clipRect = rect
}

func (surf *surface) Log(message string) {
	if surf.log == nil {
		fmt.Println(message)
		return
	}
	surf.log(message)
}

func (surf *surface) InterpolatePixel(x float32, y float32) *lang.Color {
	if x >= float32(surf.source.width) || y >= float32(surf.source.height) {
		return nil
	}
	x -= 0.5
	y -= 0.5
	if x < 0.0 || y < 0.0 {
		return nil
	}
	baseX, baseY := int(x), int(y)
	ratioX, ratioY := x-float32(baseX), y-float32(baseY)
	baseColor := surf.GetPixel(baseX, baseY)
	rightColor := surf.GetPixel(baseX+1, baseY)
	bottomColor := surf.GetPixel(baseX, baseY+1)
	return &lang.Color{
		A: baseColor.A,
		R: interpolate(baseColor.R, rightColor.R, bottomColor.R, ratioX, ratioY),
		G: interpolate(baseColor.G, rightColor.G, bottomColor.G, ratioX, ratioY),
		B: interpolate(baseColor.B, rightColor.B, bottomColor.B, ratioX, ratioY),
	}
}

func interpolate(orig lang.Number, right lang.Number, bottom lang.Number, ratioX float32, ratioY float32) lang.Number {
	dx := float64(right-orig) * float64(ratioX)
	dy := float64(bottom-orig) * float64(ratioY)
	return orig + lang.Number(math.Sqrt(dx*dx+dy*dy))
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

	return writeImage(ymg, targetFile)
}

func writeImage(ymg *ymage, writer io.Writer) error {
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

	return png.Encode(writer, img)
}
