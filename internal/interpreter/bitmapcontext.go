package interpreter

import (
	"github.com/smackem/ylang/internal/lang"
	"image"
)

// BitmapContext is the surface the ylang interpreter works on.
type BitmapContext interface {
	GetPixel(x int, y int) lang.Color
	SetPixel(x int, y int, color lang.Color)
	SourceWidth() int
	SourceHeight() int
	TargetWidth() int
	TargetHeight() int
	Convolute(x, y, width, height int, kernel []lang.Number) lang.Color
	MapRed(x, y, width, height int, kernel []lang.Number) []lang.Number
	MapGreen(x, y, width, height int, kernel []lang.Number) []lang.Number
	MapBlue(x, y, width, height int, kernel []lang.Number) []lang.Number
	MapAlpha(x, y, width, height int, kernel []lang.Number) []lang.Number
	Blt(x, y, width, height int)
	ResizeTarget(width, height int)
	Flip() int // return imageID for Recall()
	Recall(imageID int) error
	SetClipRect(rect image.Rectangle)
	ClipRect() image.Rectangle
}
